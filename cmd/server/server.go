package main

import (
	"github.com/jason-plainlog/code-exec/internal/config"
	"github.com/jason-plainlog/code-exec/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// load configuration from .env and environment variable
	config := config.GetConfig().Load()

	// initialize submission handler
	submissionHandler := &handlers.SubmissionHandler{}

	// initialize server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.BodyLimit(config.BodyLimit))

	// submission related endpoints
	e.POST("/submission", submissionHandler.Create)
	e.GET("/result/:token", submissionHandler.GetResult)

	e.Logger.Fatal(e.Start(config.Address))
}
