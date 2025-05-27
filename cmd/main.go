package main

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/mauFade/web-crawler/internal/db"
	"github.com/mauFade/web-crawler/internal/models"
	"github.com/mauFade/web-crawler/internal/utils"
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

	done := make(chan bool)
	crawlerStats := models.NewCrawlerStats("0 0\n", "0 0\n", time.Now())

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				crawlerStats.Update(crawledSet, queue, t)
				crawlerStats.Print()
			}
		}
	}()

	queue.Enqueue(seed)
	url := queue.Dequeue()
	crawledSet.Add(url)
	c := make(chan []byte)

	go utils.FetchTopLevelPage(url, c)

	content := <-c
	utils.ParseHTML(url, content, queue, crawledSet)

	for queue.GetSize() > 0 && crawledSet.GetNumber() < 5000 {
		url := queue.Dequeue()
		crawledSet.Add(url)
		go utils.FetchTopLevelPage(url, c)
		content := <-c

		if len(content) == 0 {
			continue
		}

		go utils.ParseHTML(url, content, queue, crawledSet)
	}
	ticker.Stop()
	done <- true
	dbConn.Disconnect()
	fmt.Println("\n------------------CRAWLER STATS------------------")
	fmt.Printf("Total queued: %d\n", queue.GetTotalQueued())
	fmt.Printf("To be crawled (Queue) size: %d\n", queue.GetSize())
	fmt.Printf("Crawled size: %d\n", crawledSet.GetNumber())
	crawlerStats.Print()

}
