package variable

import (
	"context"
	"fmt"
	"log"
	"os"
)

type Env string

const (
	SERVER_NAME    Env = "SERVER_NAME"
	SERVER_PORT    Env = "SERVER_PORT"
	DB_USERNAME    Env = "DB_USERNAME"
	DB_PASSWORD    Env = "DB_PASSWORD"
	DB_SCHEME      Env = "DB_SCHEME"
	DB_ADDRESS     Env = "DB_ADDRESS"
	DB_NAME        Env = "DB_NAME"
	TOKEN_SECRET   Env = "TOKEN_SECRET"
	TOKEN_DURATION Env = "TOKEN_DURATION"
	GOOGLE_API_KEY Env = "GOOGLE_API_KEY"
)

func Init() (ctx context.Context, err error) {
	ctx = context.Background()

	ctx = loadEnv(ctx, SERVER_NAME)
	ctx = loadEnv(ctx, SERVER_PORT)
	ctx = loadEnv(ctx, DB_USERNAME)
	ctx = loadEnv(ctx, DB_PASSWORD)
	ctx = loadEnv(ctx, DB_SCHEME)
	ctx = loadEnv(ctx, DB_ADDRESS)
	ctx = loadEnv(ctx, DB_NAME)
	ctx = loadEnv(ctx, TOKEN_SECRET)
	ctx = loadEnv(ctx, TOKEN_DURATION)
	ctx = loadEnv(ctx, GOOGLE_API_KEY)

	return ctx, nil
}

func loadEnv(ctx context.Context, env Env) context.Context {
	envStr := os.Getenv(string(env))
	if len(envStr) == 0 {
		log.Panicf("variable: failed to load env")
	}
	return context.WithValue(ctx, env, os.Getenv(string(env)))
}

func GetEnv(ctx context.Context, env Env) (value string) {
	return fmt.Sprintf("%v", ctx.Value(env))
}
