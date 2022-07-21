package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/talenthandongsite/server-auth/internal/util"
	"github.com/talenthandongsite/server-auth/pkg/response"
)

type SignInRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type SignInResponse struct {
	Token string `json:"token,omitempty"`
	Exp   int64  `json:"exp,omitempty"`
}

func (h *Handler) HandleSignIn(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)
	log.Println("DEBUG : in HandleSignIn")

	ctx := context.Background()

	if r.Method == http.MethodPost {
		log.Println("DEBUG : in SignIn(HandleSignIn)")
		h.PasswordSignIn(ctx, w, r)
		return
	}

	err := errors.New("method not allowed")
	log.Println(err)
	http.Error(w, err.Error(), http.StatusMethodNotAllowed)
}

func (h *Handler) PasswordSignIn(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	log.Println("DEBUG : in SignIn")

	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var signin *SignInRequest

	err = json.Unmarshal(b, &signin)
	if err != nil {
		log.Println(err)
		response.JSONError(w, err, http.StatusBadRequest)
		return
	}

	signin.Password = util.HashSHA256(signin.Password)

	user, err := h.Repo.ValidatePassword(ctx, signin.Username, signin.Password)
	if err != nil {
		log.Println(err)
		response.JSONError(w, err, http.StatusUnauthorized)
		return
	}

	token, expiration, err := h.Jwt.ForgeToken(user.ID, user.Username, user.AccessControl)
	if err != nil {
		log.Println("DEBUG : in repo ValidateUser : token forge error")
		err := errors.New("token forge error")
		response.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	data := &SignInResponse{
		Token: "Bearer " + token,
		Exp:   expiration.UnixMilli(),
	}

	databytes, err := json.Marshal(data)
	if err != nil {
		response.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	response.JSONResponse(w, databytes, http.StatusOK)
}
