package helpers

import (
	"os"
	"github.com/joho/godotenv"
	"log"
)

func EnvFileVal(key string) string {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error: loading .env values")
	}

	val := os.Getenv(key)

	if key == "MONGOURI" && val == "" {
		val = os.Getenv("MONGO_DEV_CREDS")
	}

	if val == "" {
		log.Fatal("Error: no value in .env file")
	}

	return val
}