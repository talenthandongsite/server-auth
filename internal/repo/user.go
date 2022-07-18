package repo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/talenthandongsite/server-auth/internal/util"
	"github.com/talenthandongsite/server-auth/pkg/enum/keychaintype"
	"github.com/talenthandongsite/server-auth/pkg/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DATABASE_NAME = "talent"
const USER_COLLECTION_NAME = "user"

type UserId struct{}
type KeyType struct{}
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

type KeyChainItem struct {
	Content    string `json:"content,omitempty" bson:",omitempty"`
	Secret     string `json:"secret,omitempty" bson:",omitempty"`
	Expiration int64  `json:"expiration,omitempty" bson:",omitempty"`
}

type SignIn struct {
	Username string `json:"username,omitempty" bson:",omitempty"`
	Password string `json:"password,omitempty" bson:",omitempty"`
}

type DataItem struct {
	Token string `json:"token,omitempty" bson:",omitempty"`
	Exp   int64  `json:"exp,omitempty" bson:",omitempty"`
}

type SignInResponse struct {
	Status string   `json:"status,omitempty" bson:",omitempty"`
	Data   DataItem `json:"data,omitempty" bson:",omitempty"`
}

type JWTClaims struct {
	ID            string `json:"id,omitempty" bson:"_id,omitempty"`
	Username      string `json:"username,omitempty" bson:",omitempty"`
	AccessControl string `json:"accessControl,omitempty" bson:",omitempty"`
}

type UserRepo struct {
	Coll *mongo.Collection
	Jwt  *jwt.Jwt
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

func (repo *UserRepo) ValidateUser(ctx context.Context, signin SignIn) (SignInResponse, error) {
	log.Println("DEBUG : in repo ValidateUser")

	var user User
	err := repo.Coll.FindOne(
		context.TODO(),
		bson.M{"username": signin.Username},
	).Decode(&user)

	if err != nil {
		// username에 해당되는 데이터가 없을 경우 false 반환
		log.Println("DEBUG : in repo ValidateUser : Has no matching username")
		if err == mongo.ErrNoDocuments {
			return SignInResponse{
				Status: "false",
			}, err
		}
	}

	token, expiration, err := repo.Jwt.ForgeToken(user.ID, user.Username, user.AccessControl)
	if err != nil {
		err := errors.New("token forge error")
		log.Println("DEBUG : in repo ValidateUser : token forge error")
		return SignInResponse{
			Status: "false",
		}, err
	}

	// username에 해당하는 데이터가 있고 비밀번호도 맞을 경우 리턴값 반환
	return SignInResponse{
		Status: "true",
		Data: DataItem{
			Token: "Bearer " + token,
			Exp:   expiration.UnixMilli(),
		},
	}, nil
}

func (repo *UserRepo) UpsertKeychain(ctx context.Context, keychain *KeyChainItem) (bson.M, error) {

	userId := fmt.Sprintf("%v", ctx.Value(UserId{}))
	keyType := fmt.Sprintf("%v", ctx.Value(KeyType{}))

	kct, err := keychaintype.Enum(keyType)
	if err != nil {
		return bson.M{}, err
	}

	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return bson.M{}, err
	}

	log.Println("DEBUG : in repo upsert Keychain")

	if kct == keychaintype.PASSWORD {
		keychain.Content = util.HashSHA256(keychain.Content)
		print(keychain.Content)
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"keychain." + keyType: keychain}}

	upsert := true
	after := options.After

	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	result := repo.Coll.FindOneAndUpdate(ctx, filter, update, &opt)

	if result.Err() != nil {
		return nil, result.Err()
	}

	doc := bson.M{}
	decodeErr := result.Decode(&doc)

	return doc, decodeErr
}

func (repo *UserRepo) DeleteKeychain(ctx context.Context) (bson.M, error) {

	log.Println("DEBUG : in repo Delete Keychain")

	userId := fmt.Sprintf("%v", ctx.Value(UserId{}))
	keyType := fmt.Sprintf("%v", ctx.Value(KeyType{}))

	_, err := keychaintype.Enum(keyType)
	if err != nil {
		return nil, err
	}

	objectId, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$unset": bson.M{"keychain." + keyType: ""}}

	after := options.Before

	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	result := repo.Coll.FindOneAndUpdate(ctx, filter, update, &opt)

	doc := bson.M{}
	decodeErr := result.Decode(&doc)

	return doc, decodeErr
}
