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
	e.Use(middleware.Gzip())
	e.Use(middleware.BodyLimit(config.BodyLimit))
	e.Use(middleware.Logger())

	// submission related endpoints
	e.GET("/submission/:id", submissionHandler.Get)
	e.POST("/submission", submissionHandler.Create)

	e.Logger.Fatal(e.Start(config.Address))
}
