package db

import (
	"context"
	"time"

	"github.com/mauFade/web-crawler/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

func NewMongoDB(uri, dbName, collectionName string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	collection := db.Collection(collectionName)

	return &MongoDB{
		client:     client,
		database:   db,
		collection: collection,
	}, nil
}

func (m *MongoDB) SavePage(ctx context.Context, page *models.Page) error {
	page.CreatedAt = time.Now()
	page.UpdatedAt = time.Now()

	_, err := m.collection.InsertOne(ctx, page)
	return err
}

func (m *MongoDB) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
