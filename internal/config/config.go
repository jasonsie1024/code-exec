package config

import cfg "github.com/JeremyLoy/config"

type Config struct {
	// Server Configurations
	Address    string `config:"ADDRESS" json:"-"`     // Address to listen, default to ":8000"
	ApiKey     string `config:"API_KEY" json:"-"`     // API key validation, default to ""
	BodyLimit  string `config:"BODY_LIMIT" json:"-"`  // Maximum request body size, default to "4M"
	MaxSandbox int    `config:"MAX_SANDBOX" json:"-"` // Maximum isolate sandbox

	// Storage Configurations
	SubmissionBucket string `config:"SUBMISSION_BUCKET" json:"-"` // Bucket to store submission
	TaskBucket       string `config:"TASK_BUCKET" json:"-"`       // Bucket to store task

	// Submission Configurations
	MaxTask     int     `config:"MAX_TASK" json:"max_task"`         // Maximum task amount per submission
	MaxTime     float32 `config:"MAX_TIME" json:"max_time"`         // Maximum time limit
	MaxMemory   int     `config:"MAX_MEMORY" json:"max_memory"`     // Maximum memory limit
	MaxProcess  int     `config:"MAX_PROCESS" json:"max_process"`   // Max process / thread limit
	MaxFilesize int     `config:"MAX_FILESIZE" json:"max_filesize"` // Max filesize limit
}

// The default values of the config.
var config Config = Config{
	Address:    ":8080",
	ApiKey:     "",
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
	if err := cfg.FromOptional(".env").FromEnv().To(c); err != nil {
		panic(err)
	}
	return c
}

// Get the config from config package, if not loaded yet, need to call Load method.
func GetConfig() *Config {
	return &config
}
