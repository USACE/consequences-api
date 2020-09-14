package main

import (
	"log"
	"net/http"

	"github.com/apex/gateway"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"

	"consequences-api/api/handlers"
	"consequences-api/api/middleware"
)

// Config holds all runtime configuration provided via environment variables
type Config struct {
	LambdaContext bool
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

	e := echo.New()
	e.Use(
		middleware.CORS,
		middleware.GZIP,
	)

	// Public Routes
	public := e.Group("")

	// Public Routes
	// NOTE: ALL GET REQUESTS ARE ALLOWED WITHOUT AUTHENTICATION USING JWTConfig Skipper. See appconfig/jwt.go
	public.GET("consequences/test", handlers.GetSomething())

	if cfg.LambdaContext {
		log.Print("starting server; Running On AWS LAMBDA")
		log.Fatal(gateway.ListenAndServe("localhost:3030", e))
	} else {
		log.Print("starting server")
		log.Fatal(http.ListenAndServe("localhost:3030", e))
	}
}
