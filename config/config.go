package config

import (
	"log"
	"os"
)

// Config struct contains init environments
// variables and generated passwordHash
type Config struct {
	HTTP_PORT   string
	HOTELS_PATH string
}

// NewConfig function returns inited server configuration
func NewConfig() *Config {

	HTTP_PORT := getEnv("HTTP_PORT")
	HOTELS_PATH := getEnv("HOTELS_PATH")

	return &Config{
		HTTP_PORT,
		HOTELS_PATH,
	}
}

// getEnv wrapper function to get a value from environment
func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("Env variable %s does not exist!\n", key)
		return ""
	}
	return value
}
