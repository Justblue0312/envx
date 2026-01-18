package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/justblue0312/envx"
)

type DatabaseConfig struct {
	Host     string `envx:"HOST"`
	Port     int    `envx:"PORT"`
	Password string `envx:"PASSWORD"`
}

type Config struct {
	AppName  string         `envx:"APP_NAME"`
	Database DatabaseConfig `envx:"DB" nested:"true"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Debug: Check if env vars are loaded
	fmt.Printf("FIN_APP_NAME: %s\n", os.Getenv("FIN_APP_NAME"))
	fmt.Printf("FIN_DB_HOST: %s\n", os.Getenv("FIN_DB_HOST"))
	fmt.Printf("FIN_DB_PORT: %s\n", os.Getenv("FIN_DB_PORT"))
	fmt.Printf("FIN_DB_PASSWORD: %s\n", os.Getenv("FIN_DB_PASSWORD"))

	config := &Config{}

	if err := envx.Process("FIN", config); err != nil {
		log.Fatalf("Failed to process config: %v", err)
	}

	fmt.Printf("Result: %+v\n", config)
}
