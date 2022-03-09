package main

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
	"github.com/jason-plainlog/code-exec/internal/config"
	"github.com/jason-plainlog/code-exec/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// load language.json files and config
	config.GetLanguages()
	config := config.GetConfig().Load()

	// setup storage client
	storageClient, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// initialize server & middlewares
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.BodyLimit(config.BodyLimit))

	// initizlize handlers & register routes
	infoHandler := handlers.InfoHandler{}
	infoHandler.RegisterRoutes(e)

	submissionHandler := handlers.SubmissionHandler{
		Storage: storageClient,
	}
	submissionHandler.RegisterRoutes(e)

	// start server
	e.Logger.Fatal(e.Start(config.Address))
}
