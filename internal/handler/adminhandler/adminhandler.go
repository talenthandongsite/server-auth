package adminhandler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/talenthandongsite/server-auth/internal/repo"
	"github.com/talenthandongsite/server-auth/pkg/jwt"
	"github.com/talenthandongsite/server-auth/pkg/response"
)

const RUNNER string = "adminhandler"

type AdminHandler struct {
	Repo *repo.UserRepo
	Jwt  *jwt.Jwt
}

func Init(repo *repo.UserRepo, jwt *jwt.Jwt) *AdminHandler {
	return &AdminHandler{
		Repo: repo,
		Jwt:  jwt,
	}
}

// HandleAdmin function is entry point of endpoint '/admin'. This function will act as url parser of different endpoints under '/admin'
// For example, there are some endpoints under it:
// - /admin/user: administrative api endpoint of user.
// - /admin/user/:userid/keychain/:keytype: administrative api endpoint of each user's keychain
// Job of this function is parse request and pass it to appropirate handler
func (h *AdminHandler) HandleAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	log.Println(r.Method, r.URL.Path)

	slice := strings.Split(r.URL.Path, "/")
	log.Println(slice)

	if len(slice) < 3 {
		err := errors.New(RUNNER + ": not found")
		response.JSONError(w, err, http.StatusNotFound)
		return
	}

	userId := slice[3]
	ctx = context.WithValue(ctx, repo.UserId{}, userId) //http 요청이 끝날때까지 값을 가지고있음

	if slice[2] != "user" {
		err := errors.New(RUNNER + ": not found")
		response.JSONError(w, err, http.StatusNotFound)
		return
	}

	if len(slice) < 5 {
		h.HandleUser(ctx, w, r)
		return
	}

	keyType := slice[5]
	ctx = context.WithValue(ctx, repo.KeyType{}, keyType)

	if len(slice) != 6 || slice[4] != "keychain" {
		err := errors.New(RUNNER + ": not found")
		response.JSONError(w, err, http.StatusNotFound)
		return
	}

	h.HandleKeychain(ctx, w, r)
}
