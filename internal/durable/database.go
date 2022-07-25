package durable

import (
	"context"
	"fmt"
	"log"

	"github.com/talenthandongsite/server-auth/internal/variable"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitDBClient(ctx context.Context) (*mongo.Client, error) {

	dbUsername := variable.GetEnv(ctx, variable.DB_USERNAME)
	dbPassword := variable.GetEnv(ctx, variable.DB_PASSWORD)
	dbScheme := variable.GetEnv(ctx, variable.DB_SCHEME)
	dbAddress := variable.GetEnv(ctx, variable.DB_ADDRESS)

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
