package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/talenthandongsite/server-auth/internal/repo"
	"github.com/talenthandongsite/server-auth/internal/util"
)

func (h *Handler) HandleSignIn(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)
	log.Println("DEBUG : in HandleSignIn")

	ctx := context.Background()

	if r.Method == http.MethodPost {
		log.Println("DEBUG : in SignIn(HandleSignIn)")
		h.SignIn(ctx, w, r)
		return
	}

	err := errors.New("method not allowed")
	log.Println(err)
	http.Error(w, err.Error(), http.StatusMethodNotAllowed)
}

func (h *Handler) SignIn(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	log.Println("DEBUG : in SignIn")

	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	signin := repo.SignIn{}

	err = json.Unmarshal(b, &signin)
	if err != nil {
		log.Println(err)
		log.Println("error in Unmarshalling sign in json body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	signin.Password = util.HashSHA256(signin.Password)

	validation, err := h.Repo.ValidateUser(ctx, signin)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	testBytes, err := json.Marshal(validation)
	if err != nil {
		log.Println(err)
		return
	}

	w.Write(testBytes)
}
