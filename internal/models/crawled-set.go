package models

import (
	"hash/fnv"
	"sync"
)

func hashUrl(url string) uint64 {
	hash := fnv.New64a()
	hash.Write([]byte(url))
	return hash.Sum64()
}

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

	c.Data[hashUrl(url)] = true
	c.Number++
}

func (c *CrawledSet) Contains(url string) bool {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	return c.Data[hashUrl(url)]
}

func (c *CrawledSet) GetNumber() int {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	return c.Number
}
