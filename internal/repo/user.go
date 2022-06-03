package repo

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const DATABASE_NAME = "talent"
const USER_COLLECTION_NAME = "user"

type User struct {
	ID            string         `json:"id,omitempty" bson:"_id,omitempty"`
	Username      string         `json:"username,omitempty" bson:",omitempty"`
	Password      string         `json:"password,omitempty" bson:",omitempty"`
	Email         string         `json:"email,omitempty" bson:",omitempty"`
	AccessControl string         `json:"accessControl,omitempty" bson:",omitempty"`
	Activity      []ActivityItem `bson:",omitempty"`
	KeyChain      []KeyChainItem `bson:",omitempty"`
}

type ActivityItem struct {
	Type      string `bson:",omitempty"`
	Content   string `json:"content,omitempty" bson:",omitempty"`
	TimeStamp int64  `json:"timestamp,omitempty" bson:",omitempty"`
}

type KeyChainItem struct {
	Type       string `bson:",omitempty"`
	Content    string `json:"content,omitempty" bson:",omitempty"`
	Secret     string `json:"secret,omitempty" bson:",omitempty"`
	Expiration int64  `json:"timestamp,omitempty" bson:",omitempty"`
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

func (repo *UserRepo) Update(ctx context.Context, user User, updateId string) (int, error) {
	log.Println("DEBUG : in repo Update / updateId = ", updateId)

	objectId, err := primitive.ObjectIDFromHex(updateId)
	if err != nil {
		return 0, err
	}

	update := bson.M{"$set": user}

	updateResult, err := repo.Coll.UpdateOne(
		ctx,
		bson.M{"_id": objectId},
		update,
	)
	if err != nil {
		return 0, err
	}

	log.Printf("Updated %v Document!\n", updateResult.ModifiedCount)
	updateCount := updateResult.ModifiedCount

	if updateCount == 0 {
		errorStr := "something gone wrong while updating data"
		err = errors.New(errorStr)
		return 0, err
	}
	return int(updateCount), nil
}

func (repo *UserRepo) Delete(ctx context.Context, deleteId string) (int, error) {
	log.Println("DEBUG : in repo Delete / deleteId = ", deleteId)

	objectId, err := primitive.ObjectIDFromHex(deleteId)
	if err != nil {
		return 0, err
	}

	deleteResult, err := repo.Coll.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		return 0, err
	}
	log.Printf("DeleteOne removed %v document(s)\n", deleteResult.DeletedCount)

	deletedCount := deleteResult.DeletedCount
	if deletedCount == 0 {
		errorStr := "something gone wrong while deleting data"
		err = errors.New(errorStr)
		return 0, err
	}
	return int(deletedCount), nil
}

func UpsertKeychain(ctx context.Context, id string, keychain *KeyChainItem) {

}
