package models

import (
	"fmt"
	"time"
)

type CrawlerStats struct {
	PagesPerminute        string
	CrawledRatioPerMinute string
	StartTime             time.Time
}

func NewCrawlerStats(pagesPerminute, crawledRatioPerMinute string, start time.Time) *CrawlerStats {
	return &CrawlerStats{
		PagesPerminute:        pagesPerminute,
		CrawledRatioPerMinute: crawledRatioPerMinute,
		StartTime:             start,
	}
}

func (c *CrawlerStats) Update(crawled *CrawledSet, queue *Queue, t time.Time) {
	c.PagesPerminute += fmt.Sprintf("%f %d\n", t.Sub(c.StartTime).Minutes(), crawled.GetNumber())
	c.CrawledRatioPerMinute += fmt.Sprintf("%f %f\n", t.Sub(c.StartTime).Minutes(), float64(crawled.GetNumber())/float64(queue.GetSize()))
}

func (c *CrawlerStats) Print() {
	fmt.Println("Pages crawled per minute:")
	fmt.Println(c.PagesPerminute)
	fmt.Println("Crawl to Queued Ratio per minute:")
	fmt.Println(c.CrawledRatioPerMinute)
}
