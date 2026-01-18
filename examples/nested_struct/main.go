package main

import (
	"fmt"
	"log"

	"github.com/justblue0312/envx"
)

type DatabaseConfig struct {
	Host     string `envx:"HOST"`
	Port     int    `envx:"PORT"`
	Password string `envx:"PASSWORD"`
}

type ServerConfig struct {
	Host string `envx:"HOST"`
	Port int    `envx:"PORT"`
}

type AppConfig struct {
	Name     string         `envx:"APP_NAME"`
	Debug    bool           `envx:"DEBUG"`
	Database DatabaseConfig `nested:"true"`
	Server   ServerConfig   `nested:"true"`
}

func main() {
	config := &AppConfig{}

	if err := envx.Process("", config); err != nil {
		log.Fatalf("Failed to process config: %v", err)
	}

	fmt.Printf("App Config:\n")
	fmt.Printf("  Name: %s\n", config.Name)
	fmt.Printf("  Debug: %v\n", config.Debug)
	fmt.Printf("  Database:\n")
	fmt.Printf("    Host: %s\n", config.Database.Host)
	fmt.Printf("    Port: %d\n", config.Database.Port)
	fmt.Printf("    Password: %s\n", config.Database.Password)
	fmt.Printf("  Server:\n")
	fmt.Printf("    Host: %s\n", config.Server.Host)
	fmt.Printf("    Port: %d\n", config.Server.Port)
}
