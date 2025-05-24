package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mauFade/web-crawler/internal/crawler"
	"github.com/mauFade/web-crawler/internal/db"
)

func main() {
	// Parse command line flags
	startURL := flag.String("url", "https://example.com", "Starting URL for the crawler")
	maxDepth := flag.Int("depth", 2, "Maximum crawl depth")
	mongoURI := flag.String("mongo", "mongodb://localhost:27017", "MongoDB connection URI")
	flag.Parse()

	// Create context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Received shutdown signal")
		cancel()
	}()

	// Initialize MongoDB
	mongo, err := db.NewMongoDB(*mongoURI, "webcrawler", "pages")
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongo.Close(ctx)

	// Initialize crawler
	c, err := crawler.NewCrawler(*startURL, *maxDepth)
	if err != nil {
		log.Fatalf("Failed to create crawler: %v", err)
	}

	// Start crawling
	page, err := c.Crawl(ctx, *startURL, 0)
	if err != nil {
		log.Printf("Error crawling %s: %v", *startURL, err)
		return
	}

	if page != nil {
		if err := mongo.SavePage(ctx, page); err != nil {
			log.Printf("Error saving page to MongoDB: %v", err)
		}
	}

	log.Println("Crawling completed successfully")
}
