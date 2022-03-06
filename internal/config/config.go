package config

import cfg "github.com/JeremyLoy/config"

type Config struct {
	// Server Configurations
	Address    string `config:"ADDRESS"`    // Address for web api service to listen, default to ":8000"
	BodyLimit  string `config:"BODY_LIMIT"` // Maximum request body size, default to "4M"
	MaxSandbox int    `config:"MAX_SANDBOX"`

	// Submission Configurations
	MaxTask     int     `config:"MAX_TASK"`
	MaxTime     float32 `config:"MAX_TIME"`
	MaxMemory   int     `config:"MAX_MEMORY"`
	MaxProcess  int     `config:"MAX_PROCESS"`
	MaxFilesize int     `config:"MAX_FILESIZE"`

	// Storage Configurations
	Bucket string `config:"BUCKET"`
}

// The default values of the config.
var config Config = Config{
	Address:    ":8000",
	BodyLimit:  "4M",
	MaxSandbox: 1000,

	MaxTask:     32,
	MaxTime:     15.0,
	MaxMemory:   256 * 1024,
	MaxProcess:  16,
	MaxFilesize: 4096,
}

// Load config from .env and environment variable, only need to call once in one execution.
func (c *Config) Load() *Config {
	if err := cfg.From(".env").FromEnv().To(c); err != nil {
		panic(err)
	}
	return c
}

// Get the config from config package, if not loaded yet, need to call Load method.
func GetConfig() *Config {
	return &config
}
