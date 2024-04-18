package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/kiennyo/syncwatch-be/internal/config"
	"github.com/kiennyo/syncwatch-be/internal/db"
	"github.com/kiennyo/syncwatch-be/internal/domain/users"
	"github.com/kiennyo/syncwatch-be/internal/http"
	"github.com/kiennyo/syncwatch-be/internal/mail"
	"github.com/kiennyo/syncwatch-be/internal/security"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	postgres, err := db.New(ctx, cfg.DB)
	if err != nil {
		slog.Error("Failed to connect to db", "reason", err.Error()) // Fatal
		return
	}

	mailer := mail.New(cfg.SMTP)
	tokens := security.NewTokenFactory(cfg.Security)

	// users module setup
	userRepo := users.NewRepository(postgres)
	userService := users.NewService(userRepo, tokens, mailer)
	usersHandler := users.NewHandler(userService)

	server := http.New(cfg.HTTP, tokens).
		AddRoutes("/users", usersHandler.Handlers())

	if err = server.Serve(); err != nil {
		slog.Error("Failed to start server", "reason", err.Error()) // Fatal
	}
}
