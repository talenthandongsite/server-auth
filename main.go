package main

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/talenthandongsite/server-auth/internal/durable"
	"github.com/talenthandongsite/server-auth/internal/handler"
	"github.com/talenthandongsite/server-auth/internal/repo"
	"github.com/talenthandongsite/server-auth/internal/variable"
	"github.com/talenthandongsite/server-auth/pkg/jwt"
)

const RUNNER string = "main"

func main() {

	serverStartUpCtx, err := variable.Init()
	if err != nil {
		log.Fatalln(err)
	}

	serverName := variable.GetEnv(serverStartUpCtx, variable.SERVER_NAME)
	log.Printf("%s: starting %s\n", RUNNER, serverName)

	jwtService, err := jwtInit(serverStartUpCtx)
	if err != nil {
		log.Fatalln(err)
	}
	dbclient, err := durable.InitDBClient(serverStartUpCtx)
	if err != nil {
		log.Fatalln(err)
	}

	repository := repo.InitUserRepo(serverStartUpCtx, dbclient)
	handler := handler.Init(serverStartUpCtx, repository, jwtService)

	err = handler.StartServer()
	if err != nil {
		log.Fatalln(err)
	}
}

func jwtInit(ctx context.Context) (jwtService *jwt.Jwt, err error) {
	sstr := variable.GetEnv(ctx, variable.TOKEN_SECRET)
	if len(sstr) == 0 {
		err := errors.New("main: " + (string)(variable.TOKEN_SECRET) + " didn't loaded properly")
		return nil, err
	}

	dstr := variable.GetEnv(ctx, variable.TOKEN_DURATION)
	if len(dstr) == 0 {
		err := errors.New("main: " + (string)(variable.TOKEN_DURATION) + " didn't loaded properly")
		return nil, err
	}

	secret := []byte(sstr)

	dint, err := strconv.ParseInt(dstr, 10, 64)
	if err != nil {
		return nil, err
	}
	duration := time.Millisecond * time.Duration(dint)

	config := jwt.JwtConfig{
		Secret:   secret,
		Duration: duration,
	}
	return jwt.Init(config)
}
