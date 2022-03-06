package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jason-plainlog/code-exec/internal/config"
	"github.com/jason-plainlog/code-exec/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"cloud.google.com/go/storage"
)

func main() {
	// load configuration from .env and environment variable
	config := config.GetConfig().Load()

	fmt.Println(os.Environ())

	// initialize cloud storage client
	storageClient, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create cloud storage client: %v", err)
	}

	// initialize submission handler
	submissionHandler := &handlers.SubmissionHandler{
		StorageClient: storageClient,
	}

	// initialize server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.BodyLimit(config.BodyLimit))

	// submission related endpoints
	e.POST("/submission", submissionHandler.Create)
	e.GET("/result/:token", submissionHandler.GetResult)

	e.Logger.Fatal(e.Start(config.Address))
}
