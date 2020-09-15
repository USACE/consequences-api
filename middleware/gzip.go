package middleware

import echoMiddleware "github.com/labstack/echo/v4/middleware"

// GZIP is Configured ready-to-go GZIP middleware
var GZIP = echoMiddleware.GzipWithConfig(echoMiddleware.GzipConfig{Level: 5})