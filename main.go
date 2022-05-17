package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const PORT string = "8080"
const DATABASE_NAME = "talent"
const USER_COLLECTION_NAME = "user"
const DB_URI = "mongodb+srv://talenthandongdev:nasdaq20000!@cluster0.rrf0l.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"

type User struct {
	ID            string    `json:"id,omitempty"`
	EternalID     string    `json:"externalid,omitempty"`
	UserName      string    `json:"username,omitempty"`
	PassPhrase    string    `json:"passphrase,omitempty"`
	AccessControl string    `json:"accesscontrol,omitempty" "default:"PENDING"`
	Created       time.Time `json:"created,omitempty"`
	Updated       string    `json:"updated,omitempty"`
	LastAccess    string    `json:"lastaccess,omitempty"`
	AdminNote     string    `json:"adminnote,omitempty"`
}

func main() {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(DB_URI))
	if err != nil {
		log.Println("CON ERR")
		return
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println("PING ERR")
		return
	}

	repo := InitUserRepo(client)
	handler := &UserHandler{
		Repo: repo,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Handle)
	http.ListenAndServe(":"+PORT, mux)
}

func (h *UserHandler) Handle(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)

	w.Header().Set("content-type", "application/json")

	if r.Method == http.MethodGet {
		h.Read(w, r)
		return
	}

	if r.Method == http.MethodPost {
		h.Create(w, r)
		return
	}

	err := errors.New("method not allowed")
	log.Println(err)
	http.Error(w, err.Error(), http.StatusMethodNotAllowed)
}

type UserHandler struct {
	Repo *UserRepo
}

func (h *UserHandler) Read(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	user, err := h.Repo.Read(ctx)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal struct to JSON
	jsonResp, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write to resp
	w.Write(jsonResp)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := User{AccessControl: "PENDING"}
	err = json.Unmarshal(b, &user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.Repo.Create(ctx, user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(id))
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
