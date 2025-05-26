package models

import "sync"

type Queue struct {
	TotalQueued int
	Number      int
	Elements    []string
	mutex       sync.RWMutex
}

func NewQueue(totalQueued, number int, elements []string) *Queue {
	return &Queue{
		TotalQueued: totalQueued,
		Number:      number,
		Elements:    elements,
	}
}

func (q *Queue) Enqueue(url string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.Elements = append(q.Elements, url)
	q.Number++
	q.TotalQueued++
}

func (q *Queue) Dequeue() string {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.Number == 0 {
		return ""
	}

	element := q.Elements[0]
	q.Elements = q.Elements[1:]
	q.Number--

	return element
}

func (q *Queue) GetSize() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return q.Number
}
