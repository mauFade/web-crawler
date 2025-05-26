package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/mauFade/web-crawler/internal/db"
)

func main() {
	webAccess := true

	if godotenv.Load() != nil {
		fmt.Println("Error loading .env file")
		webAccess = false
	}

	dbConn := db.NewDatabaseConnection(webAccess, "", nil, nil)
	dbConn.Connect()

	log.Println("Crawling completed successfully")
}
