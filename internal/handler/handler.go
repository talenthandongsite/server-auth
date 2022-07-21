package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/talenthandongsite/server-auth/internal/handler/adminhandler"
	"github.com/talenthandongsite/server-auth/internal/repo"
	"github.com/talenthandongsite/server-auth/pkg/jwt"
	"github.com/talenthandongsite/server-auth/pkg/variable"
)

type Handler struct {
	Repo *repo.UserRepo
	Jwt  *jwt.Jwt
	ctx  context.Context
	mux  *http.ServeMux
}

func Init(ctx context.Context, repo *repo.UserRepo, jwt *jwt.Jwt) *Handler {
	// create handler struct
	var handler *Handler = &Handler{
		Repo: repo,
		Jwt:  jwt,
		ctx:  ctx,
	}

	// create mux
	mux := http.NewServeMux()

	// create nested servers and handlers
	app := http.FileServer(http.Dir("web"))
	assets := http.FileServer(http.Dir("assets"))
	adminHandler := adminhandler.Init(repo, jwt)

	// register servers and handlers to endpoint
	mux.HandleFunc("/", handler.healthcheck)
	mux.HandleFunc("/signin", handler.HandleSignIn)
	mux.HandleFunc("/verify", handler.HandleVerify)
	mux.HandleFunc("/admin/", adminHandler.HandleAdmin)

	mux.Handle("/app", http.StripPrefix("/app", app))
	mux.Handle("/app/", http.StripPrefix("/app/", app))
	mux.Handle("/assets", http.StripPrefix("/assets", assets))
	mux.Handle("/assets/", http.StripPrefix("/assets/", assets))

	// save mux to handler struct
	handler.mux = mux

	// return handler
	return handler
}

func (h *Handler) healthcheck(w http.ResponseWriter, r *http.Request) {
	serverName := variable.GetEnv(h.ctx, variable.SERVER_NAME)
	message := "Hello from " + serverName
	w.Write([]byte(message))
}

func (h *Handler) StartServer() error {
	serverName := variable.GetEnv(h.ctx, variable.SERVER_NAME)
	serverPort := variable.GetEnv(h.ctx, variable.SERVER_PORT)

	log.Printf("%s has started on port %s\n", serverName, serverPort)
	return http.ListenAndServe(":"+serverPort, h.mux)
}
