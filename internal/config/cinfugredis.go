package config

import (
	"os"
)

type configRedis struct {
	HOST     string
	PORT     string
	PASSWORD string
	DB       string
}

func NewCfgRedis() *configRedis {

	return &configRedis{
		HOST:     os.Getenv("REDIS_HOST"),
		PORT:     os.Getenv("REDIS_PORT"),
		PASSWORD: os.Getenv("REDIS_PASSWORD"),
		DB:       os.Getenv("REDIS_DB"),
	}
}
