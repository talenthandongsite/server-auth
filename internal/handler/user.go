package handler

import (
	"context"
	"crypto"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/talenthandongsite/server-auth/internal/repo"
)

type UserHandler struct {
	Repo *repo.UserRepo
}

func InitUserHandler(repo *repo.UserRepo) *UserHandler {
	return &UserHandler{
		Repo: repo,
	}
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

func (h *UserHandler) Read(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	user, err := h.Repo.Read(ctx)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := range user {
		user[i].PassPhrase = ""
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

	nowUnixMilli := time.Now().UnixMilli()
	user := repo.User{
		AccessControl: "PENDING",
		Created:       nowUnixMilli,
		Updated:       nowUnixMilli,
		LastAccess:    nowUnixMilli,
	}
	err = json.Unmarshal(b, &user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash := crypto.SHA256.New()
	hash.Write([]byte(user.PassPhrase))
	digest := hash.Sum(nil)
	user.PassPhrase = hex.EncodeToString(digest)

	id, err := h.Repo.Create(ctx, user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(id))
}
