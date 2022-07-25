package repo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/talenthandongsite/server-auth/internal/variable"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const USER_COLLECTION_NAME = "user"

type UserId struct{}

type User struct {
	ID            string                  `json:"id,omitempty" bson:"_id,omitempty"`
	Username      string                  `json:"username,omitempty" bson:",omitempty"`
	Email         string                  `json:"email,omitempty" bson:",omitempty"`
	AccessControl string                  `json:"accessControl,omitempty" bson:",omitempty"`
	Activity      []ActivityItem          `bson:",omitempty"`
	KeyChain      map[string]KeyChainItem `bson:",omitempty"`
}

type ActivityItem struct {
	Type      string `bson:",omitempty"`
	Content   string `json:"content,omitempty" bson:",omitempty"`
	TimeStamp int64  `json:"timestamp,omitempty" bson:",omitempty"`
}

type UserRepo struct {
	Client mongo.Client
	Coll   *mongo.Collection
	ctx    context.Context
}

func InitUserRepo(ctx context.Context, client *mongo.Client) *UserRepo {
	databaseName := variable.GetEnv(ctx, variable.DB_NAME)
	database := client.Database(databaseName)
	userCollection := database.Collection(USER_COLLECTION_NAME)

	return &UserRepo{
		Client: *client,
		Coll:   userCollection,
		ctx:    ctx,
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

func (repo *UserRepo) Read(ctx context.Context, sort string, limit int64, offset int64, id string) ([]User, error) {
	var user []User
	filter := bson.M{}
	opts := options.Find()

	if len(id) > 0 {
		// id가 있는 경우에 filter에 추가
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		filter = bson.M{"_id": bson.M{"$eq": objID}}
	}

	opts.SetSkip(offset)

	opts.SetLimit(limit)

	if len(sort) > 0 {
		// TODO : 값 에러 체킹
		// 리팩토링 : 함수로 빼기
		split := strings.Split(sort, "_")
		log.Println(split[0])
		log.Println(split[1])
		key1 := split[0]
		value := 1
		if split[1] == "asc" {
			value = -1
		} else {
			value = 1
		}
		opts.SetSort(bson.M{key1: value})
	}

	cursor, err := repo.Coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &user); err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepo) Update(ctx context.Context, user User) (int, error) {

	userId := fmt.Sprintf("%v", ctx.Value(UserId{}))

	objectId, err := primitive.ObjectIDFromHex(userId)
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

func (repo *UserRepo) Delete(ctx context.Context) (int, error) {
	userId := fmt.Sprintf("%v", ctx.Value(UserId{}))

	objectId, err := primitive.ObjectIDFromHex(userId)
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
