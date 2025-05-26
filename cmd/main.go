package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/mauFade/web-crawler/internal/db"
	"github.com/mauFade/web-crawler/internal/models"
)

func main() {
	webAccess := true

	if godotenv.Load() != nil {
		fmt.Println("Error loading .env file")
		webAccess = false
	}

	dbConn := db.NewDatabaseConnection(webAccess, "", nil, nil)
	dbConn.Connect()

	crawledSet := models.NewCrawledSet(make(map[uint64]bool))
	seed := "https://www.cc.gatech.edu/"
	queue := models.NewQueue(0, 0, make([]string, 0))

	ticker := time.NewTicker(1 * time.Minute)

	log.Println("Crawling completed successfully")
}
