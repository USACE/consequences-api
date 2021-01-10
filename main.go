package main

import (
	"fmt"
	"log"
	"net/http"

	asyncer "github.com/USACE/go-simple-asyncer/asyncer"
	"github.com/apex/gateway"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/USACE/consequences-api/handlers"
	"github.com/USACE/consequences-api/middleware"
)

// Config holds all runtime configuration provided via environment variables
type Config struct {
	SkipJWT             bool
	LambdaContext       bool
	DBUser              string
	DBPass              string
	DBName              string
	DBHost              string
	DBSSLMode           string
	AsyncEngine         string `envconfig:"ASYNC_ENGINE"`
	AsyncEngineSNSTopic string `envconfig:"ASYNC_ENGINE_SNS_TOPIC"`
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

	db := Connection(
		fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s sslmode=%s binary_parameters=yes",
			cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBHost, cfg.DBSSLMode,
		),
	)

	// acquisitionAsyncer defines async engine used to package DSS files for download
	computeAsyncer, err := asyncer.NewAsyncer(
		asyncer.Config{Engine: cfg.AsyncEngine, Topic: cfg.AsyncEngineSNSTopic},
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	e := echo.New()
	e.Use(
		middleware.CORS,
		middleware.GZIP,
	)

	// Public Routes
	public := e.Group("")

	// Private Routes
	private := e.Group("")
	if cfg.SkipJWT == true {
		private.Use(middleware.MockIsLoggedIn)
	} else {
		private.Use(middleware.JWT, middleware.IsLoggedIn)
	}

	// Public Routes
	// NOTE: ALL GET REQUESTS ARE ALLOWED WITHOUT AUTHENTICATION USING JWTConfig Skipper. See appconfig/jwt.go

	// Events
	public.GET("consequences/events", handlers.ListEvents(db))
	private.POST("consequences/events", handlers.CreateEvent(db))
	private.DELETE("consequences/events/:event_id", handlers.DeleteEvent(db))

	// Computes

	// public.GET("consequences/computes", handlers.ListComputes(db))
	public.GET("consequences/computes/:compute_id", handlers.GetCompute(db))
	// public.GET("consequences/computes/:compute_id/result", handlers.GetComputeResult(db))
	private.POST("consequences/computes/bbox", handlers.RunConsequencesByBoundingBox()) //have the bbox
	private.POST("consequences/computes/fips/:fips_code/:event_id", handlers.RunConsequencesByFips(db))
	private.POST("consequences/computes/ag/xy/:year/:x/:y/:arrivaltime/:duration", handlers.RunAgConsequencesByXY())//shouldnt this be a get?

	public.GET("consequences/endpoints", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string][]string{
			"computes": []string{"events", "bbox", "fips"},
		})
	})

	if cfg.LambdaContext {
		log.Print("starting server; Running On AWS LAMBDA")
		log.Fatal(gateway.ListenAndServe("localhost:3030", e))
	} else {
		log.Print("starting server")
		log.Fatal(http.ListenAndServe("localhost:3030", e))
	}
}
