package models

import (
	"sync"

	"github.com/mauFade/web-crawler/internal/utils"
)

type CrawledSet struct {
	Data   map[uint64]bool
	Number int
	Mutex  sync.RWMutex
}

func NewCrawledSet(data map[uint64]bool) *CrawledSet {
	return &CrawledSet{
		Data:   data,
		Number: 0,
	}
}

func (c *CrawledSet) Add(url string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	c.Data[utils.HashUrl(url)] = true
	c.Number++
}

func (c *CrawledSet) Contains(url string) bool {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	return c.Data[utils.HashUrl(url)]
}

func (c *CrawledSet) GetNumber() int {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	return c.Number
}
