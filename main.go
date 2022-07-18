package main

import (
	"log"

	"github.com/talenthandongsite/server-auth/internal/durable"
	"github.com/talenthandongsite/server-auth/internal/handler"
	"github.com/talenthandongsite/server-auth/internal/repo"
	"github.com/talenthandongsite/server-auth/pkg/jwt"
	"github.com/talenthandongsite/server-auth/pkg/variable"
)

const RUNNER string = "main"

func main() {

	serverStartUpCtx, err := variable.Init()
	if err != nil {
		log.Fatalln(err)
	}

	serverName := variable.GetEnv(serverStartUpCtx, variable.SERVER_NAME)
	log.Printf("%s: starting %s\n", RUNNER, serverName)

	jwtService, err := jwt.Init(serverStartUpCtx)
	if err != nil {
		log.Fatalln(err)
	}
	dbclient, err := durable.InitDBClient(serverStartUpCtx)
	if err != nil {
		log.Fatalln(err)
	}

	repository := repo.InitUserRepo(dbclient)
	handler := handler.Init(serverStartUpCtx, repository, jwtService)

	err = handler.StartServer()
	if err != nil {
		log.Fatalln(err)
	}
}
