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
	ID            string `json:"id,omitempty" bson:"_id,omitempty"`
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Email         string `json:"email,omitempty"`
	AccessControl Type1
	Activity      []ActivityItem `bson:"inline"`
	KeyChain      []KeyChainItem `bson:"inline"` // `json:",inline" bson:",inline"`
}

type Type1 string

const (
	CREATED         Type1 = "CREATED"
	UPDATED               = "UPDATED"
	ADMIN_NOTE            = "ADMIN_NOTE"
	SIGN_IN               = "SIGN_IN"
	KEYCHAIN_UPSERT       = "KEYCHAIN_UPSERT"
	KEYCHAIN_DELETE       = "KEYCHAIN_DELETE"
)

func (t Type1) String() string {
	types := [...]string{"CREATED", "UPDATED", "ADMIN_NOTE", "SIGN_IN", "KEYCHAIN_UPSERT", "KEYCHAIN_DELETE"}

	x := string(t)
	for _, v := range types {
		if v == x {
			return x
		}
	}

	return ""
}

type AccessControl string

const (
	MASTER  AccessControl = "MASTER"
	SYSTEM                = "SYSTEM"
	ADMIN                 = "ADMIN"
	MEMBER                = "MEMBER"
	PENDING               = "PENDING"
	BANNED                = "BANNED"
)

func (a AccessControl) String() string {
	accesscontrol := [...]string{"MASTER", "SYSTEM", "ADMIN", "MEMBER", "PENDING", "BANNED"}

	x := string(a)
	for _, v := range accesscontrol {
		if v == x {
			return x
		}
	}

	return ""
}

type Type2 string

const (
	PASSWORD Type2 = "PASSWORD"
	KAKAO          = "KAKAO"
)

func (t Type2) String() string {
	types := [...]string{"PASSWORD", "KAKAO"}

	x := string(t)
	for _, v := range types {
		if v == x {
			return x
		}
	}

	return ""
}

type ActivityItem struct {
	Type1     `bson:",inline"`
	Content   string `json:"content,omitempty"`
	TimeStamp int64  `json:"timestamp,omitempty"`
}

type KeyChainItem struct {
	Type2      `bson:",inline"`
	Content    string `json:"content,omitempty"`
	Secret     string `json:"secret,omitempty"`
	Expiration int64  `json:"timestamp,omitempty"`
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

	update := bson.M{
		"$set": bson.M{
			"username":      user.Username,
			"passphrase":    user.Password,
			"accesscontrol": user.AccessControl,
			// "updated":       user.Updated,
			// "adminnote":     user.AdminNote,
		},
	}

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
