package main

import (
	"fmt"
	"os"

	"github.com/antsrp/banner_service/pkg/jwt"
	"github.com/antsrp/banner_service/pkg/logger"
	"github.com/antsrp/banner_service/pkg/logger/slog"
)

func main() {
	key := []byte(`some smart key`)
	var logger logger.Logger = slog.NewTextLogger(os.Stdout, slog.WithDebugLevel())
	js := jwt.NewJwtService(key, jwt.WithMethodHS256)

	token, err := js.NewToken(map[string]any{
		`is_admin`: true,
		`username`: `ant`,
	})

	if err != nil {
		logger.Error(fmt.Sprintf("error while create token: %v", err.Error()))
		return
	}

	fmt.Printf("token: %s\n", token)

	if err := js.Parse(token); err != nil {
		logger.Error(fmt.Sprintf("error while parse token: %v", err.Error()))
		return
	}
}
