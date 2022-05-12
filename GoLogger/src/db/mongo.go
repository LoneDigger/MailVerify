package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"me.log/src/config"
)

type Box struct {
	Type    string      `json:"type" bson:"type"`
	GetTime time.Time   `json:"get_time" bson:"get_time"`
	Data    interface{} `json:"data" bson:"data"`
}

type MongoDB struct {
	db         *mongo.Client
	collection *mongo.Collection
}

func NewMongoDB(cfg config.Mongo) *MongoDB {
	// https://stackoverflow.com/questions/55127143
	credential := options.Credential{
		Username: cfg.Username,
		Password: cfg.Password,
	}
	clientOpts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port)).SetAuth(credential)
	client, err := mongo.Connect(context.TODO(), clientOpts)

	if err != nil {
		panic(err)
	}

	if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
		panic(err)
	}

	collection := client.Database("log").Collection("log")

	return &MongoDB{
		db:         client,
		collection: collection,
	}
}

func (m *MongoDB) Write(b Box) error {
	_, err := m.collection.InsertOne(context.TODO(), b)
	return err
}
