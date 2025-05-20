package main

import (
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DatabaseConnection struct {
	access     bool
	uri        string
	client     *mongo.Client
	collection *mongo.Collection
}

func main() {
	fmt.Println("Hello, World!")
}
