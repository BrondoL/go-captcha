package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	REDIS_HOST string
	REDIS_PORT string
	REDIS_PASS string
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(err)
	}
}

func GetEnv() Config {
	cfg := Config{
		REDIS_HOST: os.Getenv("REDIS_HOST"),
		REDIS_PORT: os.Getenv("REDIS_PORT"),
		REDIS_PASS: os.Getenv("REDIS_PASS"),
	}

	return cfg
}
