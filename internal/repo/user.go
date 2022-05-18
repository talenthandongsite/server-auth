package repo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const DATABASE_NAME = "talent"
const USER_COLLECTION_NAME = "user"

type User struct {
	ID            string `json:"id,omitempty" bson:"_id,omitempty"`
	ExternalID    string `json:"externalId,omitempty"`
	Username      string `json:"username,omitempty"`
	PassPhrase    string `json:"passPhrase,omitempty"`
	AccessControl string `json:"accessControl,omitempty"`
	Created       int64  `json:"created,omitempty"`
	Updated       int64  `json:"updated,omitempty"`
	LastAccess    int64  `json:"lastAccess,omitempty"`
	AdminNote     string `json:"adminNote,omitempty"`
}

type UserRepo struct {
	Coll *mongo.Collection
}

func InitUserRepo(client *mongo.Client) *UserRepo {
	database := client.Database(DATABASE_NAME)
	userCollection := database.Collection(USER_COLLECTION_NAME)
	return &UserRepo{
		Coll: userCollection,
	}
}

func (repo *UserRepo) Create(ctx context.Context, user User) (string, error) {

	insertResult, err := repo.Coll.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	// convert inserted id into string
	objectId, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		errorStr := "something gone wrong while getting inserted id"
		err = errors.New(errorStr)
		return "", err
	}
	return objectId.Hex(), nil
}

func (repo *UserRepo) Read(ctx context.Context) ([]User, error) {
	var user []User
	cursor, err := repo.Coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &user); err != nil {
		return nil, err
	}
	return user, nil
}
