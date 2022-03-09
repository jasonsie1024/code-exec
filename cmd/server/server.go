package main

import (
	"github.com/jason-plainlog/code-exec/internal/config"
	"github.com/jason-plainlog/code-exec/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// load language.json files and config
	config.GetLanguages()
	config := config.GetConfig().Load()

	// initialize server & middlewares
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.BodyLimit(config.BodyLimit))

	// initizlize handlers & register routes
	// info handler
	infoHandler := handlers.InfoHandler{}
	infoHandler.RegisterRoutes(e)

	// start server
	e.Logger.Fatal(e.Start(config.Address))
}
