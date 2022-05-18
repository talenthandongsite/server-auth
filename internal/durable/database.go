package durable

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const DB_URI = "mongodb+srv://talenthandongdev:nasdaq20000!@cluster0.rrf0l.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"

func GetClient(ctx context.Context) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(DB_URI))
	if err != nil {
		log.Println("CON ERR")
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println("PING ERR")
		return nil, err
	}

	return client, nil
}
