package variable

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Env string

const (
	SERVER_NAME    Env = "SERVER_NAME"
	SERVER_PORT    Env = "SERVER_PORT"
	DB_USERNAME    Env = "DB_USERNAME"
	DB_PASSWORD    Env = "DB_PASSWORD"
	DB_SCHEME      Env = "DB_SCHEME"
	DB_ADDRESS     Env = "DB_ADDRESS"
	TOKEN_SECRET   Env = "TOKEN_SECRET"
	TOKEN_DURATION Env = "TOKEN_DURATION"
)

func Init() (ctx context.Context, err error) {
	ctx = context.Background()

	err = godotenv.Load("env/local.env")
	if err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, SERVER_NAME, os.Getenv(string(SERVER_NAME)))
	ctx = context.WithValue(ctx, SERVER_PORT, os.Getenv(string(SERVER_PORT)))
	ctx = context.WithValue(ctx, DB_USERNAME, os.Getenv(string(DB_USERNAME)))
	ctx = context.WithValue(ctx, DB_PASSWORD, os.Getenv(string(DB_PASSWORD)))
	ctx = context.WithValue(ctx, DB_SCHEME, os.Getenv(string(DB_SCHEME)))
	ctx = context.WithValue(ctx, DB_ADDRESS, os.Getenv(string(DB_ADDRESS)))
	ctx = context.WithValue(ctx, TOKEN_SECRET, os.Getenv(string(TOKEN_SECRET)))
	ctx = context.WithValue(ctx, TOKEN_DURATION, os.Getenv(string(TOKEN_DURATION)))

	return ctx, nil
}

func GetEnv(ctx context.Context, env Env) (value string) {
	return fmt.Sprintf("%v", ctx.Value(env))
}
