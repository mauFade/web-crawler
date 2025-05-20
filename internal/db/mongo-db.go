package db

import (
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DatabaseConnection struct {
	access     bool
	uri        string
	client     *mongo.Client
	collection *mongo.Collection
}

func (db *DatabaseConnection) Connect() {
	if db.access {
		db.uri = os.Getenv("MONGODB_URI")

		client, err := mongo.Connect(options.Client().ApplyURI(db.uri))
		if err != nil {
			panic(err)
		}

		db.client = client
		db.collection = db.client.Database("webCrawlerArchive").Collection("webpages")
	}
}
