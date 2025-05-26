package db

import (
	"context"
	"os"

	"github.com/mauFade/web-crawler/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DatabaseConnection struct {
	access     bool
	uri        string
	client     *mongo.Client
	collection *mongo.Collection
}

func NewDatabaseConnection(access bool, uri string, client *mongo.Client, collection *mongo.Collection) *DatabaseConnection {
	return &DatabaseConnection{
		access:     access,
		uri:        uri,
		client:     client,
		collection: collection,
	}
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

func (db *DatabaseConnection) Disconnect() {
	if db.access {
		db.client.Disconnect(context.Background())
	}
}

func (db *DatabaseConnection) SaveWebpage(webpage models.Webpage) {
	if db.access {
		db.collection.InsertOne(context.Background(), webpage)
	}
}
