package durable

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// context 안에 시크릿 값들을 넣고 getClient함수에서 빼서 사용
type DbUsername struct{}
type DbPassword struct{}
type DbScheme struct{}
type DbAddress struct{}

func GetClient(ctx context.Context) (*mongo.Client, error) {

	dbUsername := fmt.Sprintf("%v", ctx.Value(DbUsername{}))
	dbPassword := fmt.Sprintf("%v", ctx.Value(DbPassword{}))
	dbScheme := fmt.Sprintf("%v", ctx.Value(DbScheme{}))
	dbAddress := fmt.Sprintf("%v", ctx.Value(DbAddress{}))

	dbUri := fmt.Sprintf("%s://%s:%s@%s", dbScheme, dbUsername, dbPassword, dbAddress)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUri))
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
