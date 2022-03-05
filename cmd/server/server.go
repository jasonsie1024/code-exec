package main

import cfg "github.com/JeremyLoy/config"

func main() {
	// loading configuration from .env and environment variable
	cfg.From(".env").FromEnv().To(&config)
}
