package main

import (
	cfg "github.com/JeremyLoy/config"
	"github.com/labstack/echo/v4"
)

func main() {
	// load configuration from .env and environment variable
	cfg.From(".env").FromEnv().To(&config)

	// initialize server
	e := echo.New()

	e.Logger.Fatal(e.Start(config.Address))
}
