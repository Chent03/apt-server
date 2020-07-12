package main

import (
	"log"
	"os"
	"strconv"
)

type DbConnection struct {
	host     string
	port     int64
	user     string
	name     string
	password string
}

func getDBConnectionInfo() *DbConnection {
	port, _ := strconv.ParseInt(os.Getenv("DB_PORT"), 0, 64)
	return &DbConnection{
		host:     os.Getenv("DB_HOST"),
		port:     port,
		user:     os.Getenv("DB_USER"),
		name:     os.Getenv("DB_NAME"),
		password: os.Getenv("DB_PASSWORD"),
	}
}

// GetPortNumber returns the port number either in the .env file or process environment
func GetPortNumber() string {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("No port number provided")
	}
	return port
}
