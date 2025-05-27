package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/mauFade/web-crawler/internal/db"
	"github.com/mauFade/web-crawler/internal/models"
	"github.com/mauFade/web-crawler/internal/utils"
	"golang.org/x/net/html"
)

func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			if len(a.Val) == 0 || !strings.HasPrefix(a.Val, "http") {
				ok = false
				href = a.Val
				return ok, href
			}
			href = a.Val
			ok = true
		}
	}
	return ok, href
}

// parseHTML parses the HTML content of a webpage and extracts the title and content
// It also enqueues the links found in the webpage and saves the webpage to the database
// if the webpage has less than 1000 links crawled
func parseHTML(currUrl string, content []byte, q *models.Queue, crawled *models.CrawledSet, db *db.DatabaseConnection) {
	z := html.NewTokenizer(bytes.NewReader(content))
	tokenCount := 0
	pageContentLength := 0
	body := false
	webpage := models.Webpage{Url: currUrl, Title: "", Content: "", CreatedAt: time.Now()}
	for {
		if z.Next() == html.ErrorToken || tokenCount > 500 {
			if crawled.GetNumber() < 1000 {
				db.SaveWebpage(webpage)

			}
			return
		}
		t := z.Token()
		if t.Type == html.StartTagToken {
			if t.Data == "body" {
				body = true
			}
			if t.Data == "javascript" || t.Data == "script" || t.Data == "style" {
				// Skip script and style tags
				z.Next()
				continue
			}
			if t.Data == "title" {
				z.Next()
				title := z.Token().Data // data disappears after z.Token() is called
				webpage.Title = title
				fmt.Printf("Count: %d | %s -> %s\n", crawled.GetNumber(), currUrl, title)
			}
			if t.Data == "a" {
				ok, href := getHref(t)
				if !ok {
					continue
				}
				if crawled.Contains(href) {
					// Already crawled
					continue
				} else {
					q.Enqueue(href)
				}
			}
		}
		if body && t.Type == html.TextToken && pageContentLength < 500 {
			webpage.Content += strings.TrimSpace(t.Data)
			pageContentLength += len(t.Data)
		}
		tokenCount++
	}
}

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
	parseHTML(url, content, queue, crawledSet, dbConn)

	for queue.GetSize() > 0 && crawledSet.GetNumber() < 5000 {
		url := queue.Dequeue()
		crawledSet.Add(url)
		go utils.FetchTopLevelPage(url, c)
		content := <-c

		if len(content) == 0 {
			continue
		}

		go parseHTML(url, content, queue, crawledSet, dbConn)
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
