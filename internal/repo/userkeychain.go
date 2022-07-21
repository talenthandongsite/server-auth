package repo

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/talenthandongsite/server-auth/internal/util"
	"github.com/talenthandongsite/server-auth/pkg/enum/keychaintype"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type KeyType struct{}

type KeyChainItem struct {
	Content    string `json:"content,omitempty" bson:",omitempty"`
	Secret     string `json:"secret,omitempty" bson:",omitempty"`
	Expiration int64  `json:"expiration,omitempty" bson:",omitempty"`
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

	// to start transaction, start session
	session, err := repo.Client.StartSession()
	if err != nil {
		return bson.M{}, err
	}

	err = session.StartTransaction()
	if err != nil {
		return bson.M{}, err
	}

	var doc bson.M
	var decodeErr error

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {

		searchResult := bson.M{}
		err := repo.Coll.FindOne(ctx, bson.M{"_id": objectId}).Decode(&searchResult)
		if err != nil {
			sc.AbortTransaction(sc)
			return err
		}

		if len(searchResult) == 0 {
			sc.AbortTransaction(sc)
			err := errors.New("no such document")
			return err
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
			sc.AbortTransaction(sc)
			return result.Err()
		}

		doc = bson.M{}
		decodeErr = result.Decode(&doc)

		err = sc.CommitTransaction(sc)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return bson.M{}, err
	}

	session.EndSession(ctx)

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

func (repo *UserRepo) ValidatePassword(ctx context.Context, username string, password string) (*User, error) {

	var user *User

	err := repo.Coll.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		// username에 해당되는 데이터가 없을 경우 false 반환
		log.Println("DEBUG : in repo ValidateUser : Has no matching username")
		return nil, err
	}

	keychain, ok := user.KeyChain["password"]
	if !ok {
		err := errors.New("this user doens't have password signin method")
		return nil, err
	}

	if keychain.Content != password {
		err := errors.New("password is wrong")
		return nil, err
	}

	// username에 해당하는 데이터가 있고 비밀번호도 맞을 경우 리턴값 반환
	return user, nil
}
