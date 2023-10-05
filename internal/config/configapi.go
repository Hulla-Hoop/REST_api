package config

import (
	"os"
)

type ConfigApi struct {
	AGEAPI    string
	NATIONAPI string
	GENDERAPI string
}

func NewCfgApi() *ConfigApi {

	/* err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	} */

	return &ConfigApi{
		AGEAPI:    os.Getenv("AGEAPI"),
		NATIONAPI: os.Getenv("NATIONAPI"),
		GENDERAPI: os.Getenv("GENDERAPI"),
	}
}
