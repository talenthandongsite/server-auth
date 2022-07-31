package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/talenthandongsite/server-auth/pkg/response"
)

// Handle Signup function. This endpoint is this
// /signup/{KAKAO}
func (h *Handler) HandleSignup(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	log.Println(r.Method, r.URL.Path)

	slice := strings.Split(r.URL.Path, "/")
	log.Println(slice)

	if len(slice) < 2 {
		err := errors.New(RUNNER + ": not found")
		response.JSONError(w, err, http.StatusNotFound)
		return
	}

	method := slice[2]
	ctx = context.WithValue(ctx, "method", method)

	if r.Method == http.MethodPost {

	}

	err := errors.New("method not allowed")
	log.Println(err)
	http.Error(w, err.Error(), http.StatusMethodNotAllowed)
}
