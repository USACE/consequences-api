package main

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"

	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/structureprovider"

	"github.com/aws/aws-sdk-go/aws"
)

// Config holds all runtime configuration provided via environment variables
type Config struct {
	AWSS3Endpoint       string `envconfig:"AWS_S3_ENDPOINT"`
	AWSS3Region         string `envconfig:"AWS_S3_REGION"`
	AWSS3DisableSSL     bool   `envconfig:"AWS_S3_DISABLE_SSL"`
	AWSS3ForcePathStyle bool   `envconfig:"AWS_S3_FORCE_PATH_STYLE"`
	AWSS3Bucket         string `envconfig:"AWS_S3_BUCKET"`
}

// Connection is a database connnection
func Connection(connStr string) *sqlx.DB {
	log.Printf("Getting database connection")
	db, err := sqlx.Open("postgres", connStr)
	if err != nil || db == nil {
		log.Fatal("Could not connect to database; ", err.Error())
	}
	return db
}

func main() {
	//  Here's what would typically be here:
	// lambda.Start(Handler)
	//
	// There were a few options on how to incorporate Echo v4 on Lambda.
	//
	// Landed here for now:
	//
	//     https://github.com/apex/gateway
	//     https://github.com/labstack/echo/issues/1195
	//
	// With this for local development:
	//     https://medium.com/a-man-with-no-server/running-go-aws-lambdas-locally-with-sls-framework-and-sam-af3d648d49cb
	//
	// This looks promising and is from awslabs, but Echo v4 support isn't quite there yet.
	// There is a pull request in progress, Re-evaluate in April 2020.
	//
	//    https://github.com/awslabs/aws-lambda-go-api-proxy
	//
	var cfg Config
	if err := envconfig.Process("consequences", &cfg); err != nil {
		log.Fatal(err.Error())
	}

	// This should probably move elsewhere
	awsConfig := aws.NewConfig().WithRegion(cfg.AWSS3Region)
	// Used for "minio" during development
	awsConfig.WithDisableSSL(cfg.AWSS3DisableSSL)
	awsConfig.WithS3ForcePathStyle(cfg.AWSS3ForcePathStyle)
	if cfg.AWSS3Endpoint != "" {
		awsConfig.WithEndpoint(cfg.AWSS3Endpoint)
	}
	/*newSession, err1 := session.NewSession(awsConfig)
	if err1 != nil {
		fmt.Println(err1)
	}
	s3c := s3.New(newSession)*/

	e := echo.New()
	e.Use(
		middleware.CORS(),
		middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}),
	)

	// Public Routes
	public := e.Group("")

	// Private Routes
	/*private := e.Group("")
	if cfg.SkipJWT == true {
		private.Use(middleware.MockIsLoggedIn)
	} else {
		private.Use(middleware.JWT, middleware.IsLoggedIn)
	}*/

	// Public Routes
	// NOTE: ALL GET REQUESTS ARE ALLOWED WITHOUT AUTHENTICATION USING JWTConfig Skipper. See appconfig/jwt.go
	public.GET("consequences", func(c echo.Context) error {
		return c.String(http.StatusOK, "consequences-api v0.0.1")
	})
	public.POST("consequences/summary/compute", func(c echo.Context) error {
		var i Compute
		if err := c.Bind(&i); err != nil {
			return c.String(http.StatusBadRequest, "Invalid Input")
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusOK)
		if !i.valid() {
			return c.String(http.StatusBadRequest, "File Path is invalid")
		}
		var sp consequences.StreamProvider
		if i.InventorySource == "" || i.InventorySource == "NSI" {
			sp = structureprovider.InitNSISP()
		}
		if len(i.InventorySource) > 3 {
			if i.InventorySource[len(i.InventorySource)-3:] == "shp" {
				sp = structureprovider.InitSHP(i.InventorySource)
			}
			if i.InventorySource[len(i.InventorySource)-4:] == "gpgk" {
				sp = structureprovider.InitGPK(i.InventorySource, "nsi")
			}
		}
		rw := consequences.InitSummaryResultsWriter(c.Response())
		//if output type is not summary or blank throw error?
		if i.OutputType == "" || i.OutputType == "Summary" {
			dfr := hazardproviders.Init(i.DepthFilePath)
			compute.StreamAbstract(dfr, sp, rw)
			return c.NoContent(http.StatusOK)
		}
		return c.String(http.StatusBadRequest, "OutputType must be blank or Summary")

	})
	public.POST("consequences/structure/compute", func(c echo.Context) error {
		var i Compute
		if err := c.Bind(&i); err != nil {
			return c.String(http.StatusBadRequest, "Invalid Input")
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusOK)
		if !i.valid() {
			return c.String(http.StatusBadRequest, "File Path is invalid")
		}
		var sp consequences.StreamProvider
		if i.InventorySource == "" || i.InventorySource == "NSI" {
			sp = structureprovider.InitNSISP()
		}
		if len(i.InventorySource) > 3 {
			if i.InventorySource[len(i.InventorySource)-3:] == "shp" {
				sp = structureprovider.InitSHP(i.InventorySource)
			}
			if i.InventorySource[len(i.InventorySource)-4:] == "gpgk" {
				sp = structureprovider.InitGPK(i.InventorySource, "nsi")
			}
		}

		var rw consequences.ResultsWriter
		rw = consequences.InitStreamingResultsWriter(c.Response())
		if i.OutputType == "Summary" {
			return c.String(http.StatusBadRequest, "Summary output type detected - please use consequences/summary/compute")
		}
		if i.OutputType == "GeoJson" {
			rw = consequences.InitGeoJsonResultsWriter(c.Response())
		}
		dfr := hazardproviders.Init(i.DepthFilePath)
		compute.StreamAbstract(dfr, sp, rw)
		return c.NoContent(http.StatusOK)
	})

	log.Print("starting server")
	log.Fatal(http.ListenAndServe(":8000", e))
}

type Compute struct {
	Name            string `json:"name"`
	DepthFilePath   string `json:"depthfilepath"`
	InventorySource string `json:"inventorysource"`
	OutputType      string `json:"outputtype"`
}

func (c Compute) valid() bool {

	//validate that the input depth file path is a tif and has vsis3?
	//validate the inventory path is NSI, a valid shp or a valid geopackage?
	return true //@TODO implement me!
}
