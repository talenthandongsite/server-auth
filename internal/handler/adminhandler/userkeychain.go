package adminhandler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/talenthandongsite/server-auth/internal/repo"
)

// TODO: implement handler about user/keychain
func (h *AdminHandler) HandleKeychain(ctx context.Context, w http.ResponseWriter, r *http.Request) {

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

func (h *AdminHandler) KeychainUpsert(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

func (h *AdminHandler) KeychainDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", "application/json")

	doc, _ := h.Repo.DeleteKeychain(ctx)

	jsonString, _ := json.Marshal(doc["keychain"])

	w.Write([]byte(jsonString))

}
