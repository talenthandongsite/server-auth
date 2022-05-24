package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/talenthandongsite/server-auth/internal/durable"
	"github.com/talenthandongsite/server-auth/internal/handler"
	"github.com/talenthandongsite/server-auth/internal/repo"
)

const PORT string = "8080"

func main() {
	ctx := context.Background()

	fmt.Println("Starting Talent Server")

	client, err := durable.GetClient(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	repo := repo.InitUserRepo(client)
	handler := handler.InitUserHandler(repo)

	mux := http.NewServeMux()

	mux.HandleFunc("/admin/user", handler.HandleCreateRead)
	mux.HandleFunc("/admin/user/", handler.HandleUpdateDelete)
	// mux.HandleFunc(("/"))

	http.ListenAndServe(":"+PORT, mux)
}
