package repo

import (
	"context"
	"errors"
	"log"

	"github.com/dgrijalva/jwt-go"
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
	Data   DataItem `json:"data,omitempty, bson:",omitempty"`
}

type JWTClaims struct {
	ID                 string `json:"id,omitempty" bson:"_id,omitempty"`
	Username           string `json:"username,omitempty" bson:",omitempty"`
	AccessControl      string `json:"accessControl,omitempty" bson:",omitempty"`
	jwt.StandardClaims        // 표준 토큰 Claims
}

type UserRepo struct {
	Coll *mongo.Collection
}

type Gotest struct {
	Username string `bson:"username"`
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

func (repo *UserRepo) ValidateUser(ctx context.Context, signin SignIn) (SignInResponse, error) {
	log.Println("DEBUG : in repo ValidateUser")

	var result bson.M
	err := repo.Coll.FindOne(
		context.TODO(),
		bson.D{{"username", signin.Username}},
	).Decode(&result)

	if err != nil {
		// username에 해당되는 데이터가 없을 경우 false 반환
		log.Println("DEBUG : in repo ValidateUser : Has no matching username")
		if err == mongo.ErrNoDocuments {
			return SignInResponse{
				Status: "false",
			}, err
		}
	}

	if result["password"] != signin.Password {
		// username에 해당되는 데이터가 있지만 password가 틀릴 경우 false 반환
		err := errors.New("Wrong password")
		log.Println("DEBUG : in repo ValidateUser : Wrong password")
		return SignInResponse{
			Status: "false",
		}, err
	}

	id := result["_id"].(primitive.ObjectID).Hex()
	username := result["username"].(string)
	accesscontrol := result["accesscontrol"].(string)

	// username에 해당하는 데이터가 있고 비밀번호도 맞을 경우 리턴값 반환
	return SignInResponse{
		Status: "true",
		Data: DataItem{
			Token: "bearing " + TestTokenBuild(id, username, accesscontrol),
			Exp:   1234456,
		},
	}, nil
}

func TestTokenBuild(id string, username string, accesscontrol string) string {
	at := JWTClaims{
		ID:            id,
		Username:      username,
		AccessControl: accesscontrol,
		// StandardClaims: jwt.StandardClaims{
		// 	ExpiresAt: jwt.At(time.Now().Add(time.Hour * 24)), // 만료시간 24시간
		// },
	}

	atoken := jwt.NewWithClaims(jwt.SigningMethodHS256, &at)
	signedAuthToken, err := atoken.SignedString([]byte("SecretCode"))

	if err != nil {
		log.Println(err)
	}

	log.Println("-------------------------------")
	log.Println("bearing " + signedAuthToken)

	return signedAuthToken
}

func UpsertKeychain(ctx context.Context, id string, keychain *KeyChainItem) {

}
