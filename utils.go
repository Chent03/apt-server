package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// GetPortNumber returns the port number either in the .env file or process environment
func GetPortNumber() string {
	if err := godotenv.Load(); err != nil {
		log.Print(".env file missing")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("No port number provided")
	}
	return port
}
