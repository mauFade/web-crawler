package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Page represents a crawled web page
type Page struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	URL       string             `bson:"url"`
	Title     string             `bson:"title"`
	Content   string             `bson:"content"`
	Links     []string           `bson:"links"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}
