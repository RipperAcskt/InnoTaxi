package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/RipperAcskt/innotaxi/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Mongo struct {
	client *mongo.Client
	cfg    *config.Config
}

func New(cfg *config.Config) (*Mongo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.GetMongoUrl()))
	if err != nil {
		return nil, fmt.Errorf("connect failed: %w", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	return &Mongo{client, cfg}, nil
}

func (m *Mongo) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := m.client.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("disconnect failed: %w", err)
	}
	return nil
}

func (m *Mongo) Write(p []byte) (n int, err error) {
	logs := m.client.Database("innotaxi").Collection("logs")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := logs.InsertOne(ctx, bson.M{
		"created": time.Now(),
		"log":     p,
	})
	fmt.Println(err, string(p), res.InsertedID)
	if err != nil {
		return
	}

	return len(p), nil
}
