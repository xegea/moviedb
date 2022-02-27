package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env    string
	Port   string
	ApiKey string
	ApiUrl string
}

func LoadConfig(env *string) (Config, error) {

	err := godotenv.Load(*env)
	if err != nil {
		log.Printf("Error loading %s file", *env)
	}

	environment := os.Getenv("ENV")

	port := os.Getenv("PORT")
	if port == "" {
		return Config{}, fmt.Errorf("PORT cannot be empty")
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return Config{}, fmt.Errorf("API_KEY cannot be empty")
	}

	apiUrl := os.Getenv("API_URL")
	if apiUrl == "" {
		return Config{}, fmt.Errorf("API_URL cannot be empty")
	}

	return Config{
		Env:    environment,
		Port:   port,
		ApiKey: apiKey,
		ApiUrl: apiUrl,
	}, nil
}
