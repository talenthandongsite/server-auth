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
	"github.com/talenthandongsite/server-auth/pkg/jwt"

	"github.com/joho/godotenv"
)

const PORT string = "8080"

func main() {
	err := godotenv.Load("env/local.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	ctx := context.Background()

	ctx = context.WithValue(ctx, durable.DbUsername{}, os.Getenv("DB_USERNAME"))
	ctx = context.WithValue(ctx, durable.DbPassword{}, os.Getenv("DB_PASSWORD"))
	ctx = context.WithValue(ctx, durable.DbScheme{}, os.Getenv("DB_SCHEME"))
	ctx = context.WithValue(ctx, durable.DbAddress{}, os.Getenv("DB_ADDRESS"))
	ctx = context.WithValue(ctx, jwt.TokenSecret{}, os.Getenv("TOKEN_SECRET"))
	ctx = context.WithValue(ctx, jwt.TokenDuration{}, os.Getenv("TOKEN_DURATION"))

	fmt.Println("Starting Talent Server")

	jwtService, err := jwt.Init(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	dbclient, err := durable.GetDBClient(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	mux := http.NewServeMux()
	repository := repo.InitUserRepo(dbclient)
	handler := handler.InitUserHandler(repository, jwtService)

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
