package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configkafka struct {
	BootstrapService string
	GroupID          string
	AutoOffsetReset  string
	Topic            string
	TopicErr         string
}

func New() *Configkafka {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Configkafka{
		BootstrapService: os.Getenv("BOOTSTRAPSERVER"),
		AutoOffsetReset:  os.Getenv("AUTOOFFSETRESET"),
		GroupID:          os.Getenv("GROUPID"),
		Topic:            os.Getenv("TOPIC"),
		TopicErr:         os.Getenv("TOPICERR"),
	}
}
