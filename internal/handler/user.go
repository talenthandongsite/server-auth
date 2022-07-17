package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/talenthandongsite/server-auth/internal/repo"
	"github.com/talenthandongsite/server-auth/internal/util"
	"github.com/talenthandongsite/server-auth/pkg/enum/accesscontrol"
	"github.com/talenthandongsite/server-auth/pkg/jwt"
)

type UserHandler struct {
	Repo *repo.UserRepo
	Jwt  *jwt.Jwt
}

func InitUserHandler(repo *repo.UserRepo, jwt *jwt.Jwt) *UserHandler {
	return &UserHandler{
		Repo: repo,
		Jwt:  jwt,
	}
}

func (h *UserHandler) HandleUser(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)
	log.Println("DEBUG : in HandleCreateRead")
	w.Header().Set("content-type", "application/json")

	slice := strings.Split(r.URL.Path, "/")

	if len(slice) == 6 {
		log.Println("DEBUG : in HandleKeychain")
		h.HandleKeychain(w, r)
		return
	}

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

func (h *UserHandler) HandleSignIn(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)
	log.Println("DEBUG : in HandleSignIn")

	w.Header().Set("content-type", "application/json")

	if r.Method == http.MethodPost {
		log.Println("DEBUG : in SignIn(HandleSignIn)")
		h.SignIn(w, r)
		return
	}

	err := errors.New("method not allowed")
	log.Println(err)
	http.Error(w, err.Error(), http.StatusMethodNotAllowed)
}

func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
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

func (h *UserHandler) Read(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	log.Println("DEBUG : in Read")

	var sort = ""
	var limit int64 = 10
	var offset int64 = 0
	var id = ""

	query, err1 := url.ParseQuery(r.URL.RawQuery)
	if err1 != nil {
		panic(err1)
	}
	log.Println(query)

	// 1. 있는지 체크
	// 있으면 val값, true
	// 없으면 nil, false
	// TODO : 값 에러 체킹
	if val, ok := query["sort"]; ok {
		sort = val[0]
	}

	if val, ok := query["limit"]; ok {
		limit, _ = strconv.ParseInt(val[0], 10, 64)
	}

	if val, ok := query["offset"]; ok {
		offset, _ = strconv.ParseInt(val[0], 10, 64)
	}

	if val, ok := query["id"]; ok {
		id = val[0]
	}

	// log.Println(sort)
	// log.Println(limit)
	// log.Println(offset)
	// log.Println(id)

	user, err := h.Repo.Read(ctx, sort, limit, offset, id)
	// TODO : 파라미터 구조체로 넘기기
	if err != nil {
		log.Println(err)
		log.Println("here")
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
		log.Println("hello")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = accesscontrol.Enum(user.AccessControl)
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

// token validation?
func (h *UserHandler) HandleAuth(w http.ResponseWriter, r *http.Request) {
	// 요청으로부터 쿠키의 토큰 가져오기
	token, ok := r.Header["Authorization"]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get the JWT string from the cookie
	tknStr := token[0]
	tokenStr := strings.Split(tknStr, " ")
	log.Println(tokenStr[1])

	// Initialize a new instance of `Claims`
	claims, err := h.Jwt.OpenToken(tokenStr[1])
	log.Println(claims, err)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// 인증 성공한 경우 Welcome message 출력
	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
}
