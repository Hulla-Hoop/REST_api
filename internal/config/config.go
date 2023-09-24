package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	BootstrapService string
	GroupID          string
	AutoOffsetReset  string
}

func New() *config {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &config{
		BootstrapService: os.Getenv("BOOTSTRAPSERVER"),
		AutoOffsetReset:  os.Getenv("AUTOOFFSETRESET"),
		GroupID:          os.Getenv("GROUPID"),
	}
}
