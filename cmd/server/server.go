package main

import (
	"context"
	"log"

	"github.com/jason-plainlog/code-exec/internal/config"
	"github.com/jason-plainlog/code-exec/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"cloud.google.com/go/storage"
)

func main() {
	// load configuration from .env and environment variable
	config := config.GetConfig().Load()

	// initialize cloud storage client
	storageClient, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create cloud storage client: %v", err)
	}

	// initialize submission handler
	submissionHandler := &handlers.SubmissionHandler{
		StorageClient: storageClient,
	}
	infoHandler := &handlers.InfoHandler{}

	// initialize server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.BodyLimit(config.BodyLimit))

	// submission related endpoints
	e.POST("/submission", submissionHandler.Create)
	e.GET("/result/:token", submissionHandler.GetResult)

	// info related endpoints
	e.GET("/config", infoHandler.ConfigInfo)

	e.Logger.Fatal(e.Start(config.Address))
}
