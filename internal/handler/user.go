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
	"strconv"
	"strings"

	// "time"

	"github.com/talenthandongsite/server-auth/internal/enum/accesscontrol"
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

func (h *UserHandler) HandleCreateRead(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)
	log.Println("DEBUG : in HandleCreateRead")

	w.Header().Set("content-type", "application/json")

	if r.Method == http.MethodGet {
		log.Println("DEBUG : in Read(HandleCreateRead)")
		h.Read(w, r)
		return
	}

	if r.Method == http.MethodPost {
		log.Println("DEBUG : in Create(HandleCreateRead)")
		h.Create(w, r)
		return
	}

	err := errors.New("method not allowed")
	log.Println(err)
	http.Error(w, err.Error(), http.StatusMethodNotAllowed)
}

func (h *UserHandler) HandleUpdateDelete(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)
	log.Println("DEBUG : in HandleUpdateDelete")

	w.Header().Set("content-type", "application/json")

	if r.Method == http.MethodPut {
		log.Println("DEBUG : in Update(HandleUpdateDelete)")
		h.Update(w, r)
		return
	}

	if r.Method == http.MethodDelete {
		log.Println("DEBUG : in Delete(HandleUpdateDelete)")
		h.Delete(w, r)
		return
	}

	err := errors.New("method not allowed")
	log.Println(err)
	http.Error(w, err.Error(), http.StatusMethodNotAllowed)
}

func (h *UserHandler) Read(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	log.Println("DEBUG : in Read")

	user, err := h.Repo.Read(ctx)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := range user {
		user[i].Password = ""
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
	log.Println("DEBUG : in Create")

	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// nowUnixMilli := time.Now().UnixMilli()
	user := repo.User{
		AccessControl: "PENDING",
		// Created:       nowUnixMilli,
		// Updated:       nowUnixMilli,
		// LastAccess:    nowUnixMilli,
	}
	err = json.Unmarshal(b, &user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = accesscontrol.Enum(user.AccessControl)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash := crypto.SHA256.New()
	hash.Write([]byte(user.Password))
	digest := hash.Sum(nil)
	user.Password = hex.EncodeToString(digest)

	id, err := h.Repo.Create(ctx, user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(id))
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	log.Println("DEBUG : in Update")

	slice := strings.Split(r.URL.Path, "/")
	updateId := slice[len(slice)-1]

	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// nowUnixMilli := time.Now().UnixMilli()
	user := repo.User{
		AccessControl: "PENDING",
		// Updated:       nowUnixMilli,
	}
	err = json.Unmarshal(b, &user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash := crypto.SHA256.New()
	hash.Write([]byte(user.Password))
	digest := hash.Sum(nil)
	user.Password = hex.EncodeToString(digest)

	count, err := h.Repo.Update(ctx, user, updateId)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updateMessage := "Updated " + strconv.Itoa(count) + " document"
	w.Write([]byte(updateMessage))
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	log.Println("DEBUG : in Delete")

	slice := strings.Split(r.URL.Path, "/")
	deleteId := slice[len(slice)-1]

	count, err := h.Repo.Delete(ctx, deleteId)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deleteMessage := "Deleted " + strconv.Itoa(count) + " document"
	w.Write([]byte(deleteMessage))
}
