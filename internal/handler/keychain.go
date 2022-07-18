package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/talenthandongsite/server-auth/internal/repo"
)

// TODO: implement handler about user/keychain

func (h *UserHandler) HandleKeychain(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	slice := strings.Split(r.URL.Path, "/")

	userId := slice[3]
	keyType := slice[5]

	ctx = context.WithValue(ctx, repo.UserId{}, userId) //http 요청이 끝날때까지 값을 가지고있음
	ctx = context.WithValue(ctx, repo.KeyType{}, keyType)

	if r.Method == http.MethodPut {
		log.Println("DEBUG : Upsert Keychain")
		h.KeychainUpsert(ctx, w, r)
		return
	}

	if r.Method == http.MethodDelete {
		log.Println("DEBUG : Delete keychain")
		h.KeychainDelete(ctx, w, r)
		return
	}

	err := errors.New("method not allowed")
	log.Println(err)
	http.Error(w, err.Error(), http.StatusMethodNotAllowed)

}

func (h *UserHandler) KeychainUpsert(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)
	log.Println("DEBUG : in Handle KeyChainUpsert")

	w.Header().Set("content-type", "application/json")

	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	keychain := repo.KeyChainItem{}

	err = json.Unmarshal(b, &keychain)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	doc, err := h.Repo.UpsertKeychain(ctx, &keychain)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonString, err := json.Marshal(doc["keychain"])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(jsonString))
}

func (h *UserHandler) KeychainDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)
	log.Println("DEBUG : in Handle KeyChainDelete")

	w.Header().Set("content-type", "application/json")

	doc, _ := h.Repo.DeleteKeychain(ctx)

	jsonString, _ := json.Marshal(doc["keychain"])

	w.Write([]byte(jsonString))

}
