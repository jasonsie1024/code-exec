package main

type Config struct {
	Address string `config:"ADDRESS"` // Address for web api service to listen, default to ":8000"
}

var config Config = Config{
	Address: ":8000",
}
