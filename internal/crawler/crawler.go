package crawler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/mauFade/web-crawler/internal/models"
	"golang.org/x/net/html"
)

type Crawler struct {
	client      *http.Client
	visited     map[string]bool
	visitedLock sync.RWMutex
	maxDepth    int
	baseURL     *url.URL
}

func NewCrawler(baseURL string, maxDepth int) (*Crawler, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &Crawler{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		visited:  make(map[string]bool),
		maxDepth: maxDepth,
		baseURL:  parsedURL,
	}, nil
}

func (c *Crawler) Crawl(ctx context.Context, url string, depth int) (*models.Page, error) {
	if depth > c.maxDepth {
		return nil, nil
	}

	c.visitedLock.RLock()
	if c.visited[url] {
		c.visitedLock.RUnlock()
		return nil, nil
	}
	c.visitedLock.RUnlock()

	c.visitedLock.Lock()
	c.visited[url] = true
	c.visitedLock.Unlock()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	page := &models.Page{
		URL:   url,
		Links: make([]string, 0),
	}

	var title string
	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "title" && n.FirstChild != nil {
				title = n.FirstChild.Data
			}
			if n.Data == "a" {
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						link, err := url.Parse(attr.Val)
						if err != nil {
							continue
						}
						absoluteURL := c.baseURL.ResolveReference(link).String()
						if strings.HasPrefix(absoluteURL, c.baseURL.Scheme+"://"+c.baseURL.Host) {
							links = append(links, absoluteURL)
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	page.Title = title
	page.Links = links

	return page, nil
}
