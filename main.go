package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/talenthandongsite/server-auth/internal/durable"
	"github.com/talenthandongsite/server-auth/internal/handler"
	"github.com/talenthandongsite/server-auth/internal/repo"

	"github.com/joho/godotenv"
)

const PORT string = "8080"

func main() {
	err := godotenv.Load("env/local.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	val := os.Getenv("DB_USERNAME")
	fmt.Println(val)

	ctx := context.Background()

	fmt.Println("Starting Talent Server")

	client, err := durable.GetClient(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	mux := http.NewServeMux()

	repository := repo.InitUserRepo(client)

	handler := handler.InitUserHandler(repository)

	mux.HandleFunc("/admin/user", handler.HandleUser)
	mux.HandleFunc("/admin/user/", handler.HandleUser)

	mux.HandleFunc("/signin", handler.HandleSignIn)
	mux.HandleFunc("/auth", handler.HandleAuth)

	app := http.FileServer(http.Dir("web"))
	assets := http.FileServer(http.Dir("assets"))

	mux.Handle("/app", http.StripPrefix("/app", app))
	mux.Handle("/app/", http.StripPrefix("/app/", app))
	mux.Handle("/assets", http.StripPrefix("/assets", assets))
	mux.Handle("/assets/", http.StripPrefix("/assets/", assets))

	http.ListenAndServe(":"+PORT, mux)
}
