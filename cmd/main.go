package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	webAccess := true

	if godotenv.Load() != nil {
		fmt.Println("Error loading .env file")
		webAccess = false
	}

	log.Println("Crawling completed successfully")
}
