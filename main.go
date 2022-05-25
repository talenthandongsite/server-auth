package main

import (
	"net/http"
)

const PORT string = "8080"

func main() {
	// ctx := context.Background()

	// client, err := durable.GetClient(ctx)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	mux := http.NewServeMux()

	// repo := repo.InitUserRepo(client)
	// handler := handler.InitUserHandler(repo)

	// mux.HandleFunc("/", handler.Handle)

	app := http.FileServer(http.Dir("web"))
	assets := http.FileServer(http.Dir("assets"))

	mux.Handle("/app", http.StripPrefix("/app", app))
	mux.Handle("/app/", http.StripPrefix("/app/", app))
	mux.Handle("/assets", http.StripPrefix("/assets", assets))
	mux.Handle("/assets/", http.StripPrefix("/assets/", assets))

	http.ListenAndServe(":"+PORT, mux)
}
