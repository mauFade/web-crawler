package utils

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/mauFade/web-crawler/internal/models"
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

func ParseHTML(currUrl string, content []byte, q *models.Queue, crawled *models.CrawledSet) {
	z := html.NewTokenizer(bytes.NewReader(content))
	tokenCount := 0
	pageContentLength := 0
	body := false
	webpage := models.Webpage{Url: currUrl, Title: "", Content: ""}
	for {
		if z.Next() == html.ErrorToken || tokenCount > 500 {
			if crawled.GetNumber() < 1000 {
				// db.SaveWebpage(webpage)
				fmt.Println(webpage)
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
