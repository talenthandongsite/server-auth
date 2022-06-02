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

	mux := http.NewServeMux()

	repository := repo.InitUserRepo(client)

	handler := handler.InitUserHandler(repository)

	mux.HandleFunc("/admin/user", handler.HandleCreateRead)
	mux.HandleFunc("/admin/user/", handler.HandleUpdateDelete)

	app := http.FileServer(http.Dir("web"))
	assets := http.FileServer(http.Dir("assets"))

	mux.Handle("/app", http.StripPrefix("/app", app))
	mux.Handle("/app/", http.StripPrefix("/app/", app))
	mux.Handle("/assets", http.StripPrefix("/assets", assets))
	mux.Handle("/assets/", http.StripPrefix("/assets/", assets))

	var aa repo.AccessControl = repo.PENDING
	fmt.Print(aa)

	var t repo.Type1 = repo.ADMIN_NOTE
	fmt.Print(t)

	// switch t {
	// case repo.MASTER:
	// 	fmt.Println("Master Picked")
	// case repo.ADMIN:
	// 	fmt.Println("Admin Picked")
	// default:
	// 	fmt.Println("Default")
	// }
	// fmt.Print(repo.MASTER.String())

	http.ListenAndServe(":"+PORT, mux)
}
