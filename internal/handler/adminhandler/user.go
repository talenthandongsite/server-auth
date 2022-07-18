package adminhandler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/talenthandongsite/server-auth/internal/repo"
	"github.com/talenthandongsite/server-auth/pkg/enum/accesscontrol"
)

func (h *AdminHandler) HandleUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)
	log.Println("DEBUG : in HandleCreateRead")
	w.Header().Set("content-type", "application/json")

	if r.Method == http.MethodGet {
		log.Println("DEBUG : in Read(HandleCreateRead)")
		h.Read(ctx, w, r)
		return
	}

	if r.Method == http.MethodPost {
		log.Println("DEBUG : in Create(HandleCreateRead)")
		h.Create(ctx, w, r)
		return
	}

	if r.Method == http.MethodPut {
		log.Println("DEBUG : in Update(HandleUpdateDelete)")
		h.Update(ctx, w, r)
		return
	}

	if r.Method == http.MethodDelete {
		log.Println("DEBUG : in Delete(HandleUpdateDelete)")
		h.Delete(ctx, w, r)
		return
	}

	err := errors.New("method not allowed")
	log.Println(err)
	http.Error(w, err.Error(), http.StatusMethodNotAllowed)
}

func (h *AdminHandler) Read(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

func (h *AdminHandler) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

func (h *AdminHandler) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	log.Println("DEBUG : in Update")

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

	count, err := h.Repo.Update(ctx, user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updateMessage := "Updated " + strconv.Itoa(count) + " document"
	w.Write([]byte(updateMessage))
}

func (h *AdminHandler) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	log.Println("DEBUG : in Delete")

	count, err := h.Repo.Delete(ctx)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deleteMessage := "Deleted " + strconv.Itoa(count) + " document"
	w.Write([]byte(deleteMessage))
}
