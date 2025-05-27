package models

import "time"

type Webpage struct {
	Url       string
	Title     string
	Content   string
	CreatedAt time.Time
}
