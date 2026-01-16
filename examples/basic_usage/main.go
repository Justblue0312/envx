package main

import (
	"fmt"
	"log"
	"os"

	"github.com/envx/envx"
)

type ServiceConfig struct {
	Name     string `envconfig:"SERVICE_NAME"`
	Port     int    `envconfig:"SERVICE_PORT"`
	Debug    bool   `envconfig:"DEBUG" default:"false"`
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
}

func main() {
	os.Setenv("APP_SERVICE_NAME", "my-service")
	os.Setenv("APP_SERVICE_PORT", "8080")
	os.Setenv("APP_DEBUG", "true")

	config := &ServiceConfig{}

	if err := envx.Process("APP", config); err != nil {
		log.Fatalf("Failed to process config: %v", err)
	}

	fmt.Printf("Service Configuration:\n")
	fmt.Printf("  Name: %s\n", config.Name)
	fmt.Printf("  Port: %d\n", config.Port)
	fmt.Printf("  Debug: %v\n", config.Debug)
	fmt.Printf("  LogLevel: %s\n", config.LogLevel)

	if err := envx.CheckDisallowed("APP", config); err != nil {
		fmt.Printf("Warning: %v\n", err)
	} else {
		fmt.Println("No disallowed environment variables found")
	}
}
