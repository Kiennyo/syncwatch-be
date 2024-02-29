package main

import (
	"log/slog"
	"sync"

	"github.com/kiennyo/syncwatch-be/internal/domain/users"
	"github.com/kiennyo/syncwatch-be/internal/infrastructure/config"
	"github.com/kiennyo/syncwatch-be/internal/infrastructure/http"
	"github.com/kiennyo/syncwatch-be/internal/infrastructure/log"
	"github.com/kiennyo/syncwatch-be/internal/infrastructure/security"
)

var wg sync.WaitGroup

func main() {
	log.Init()
	cfg := config.Load()

	tokens := security.NewTokens(cfg.Security.JWTSecret)

	authHandler := users.NewHandler()

	server := http.New(&wg, cfg.HTTP, tokens).
		AddRoutes("/users", authHandler.Handlers())

	if err := server.Serve(); err != nil {
		slog.Error("Failed to start server", "reason", err.Error()) // Fatal
	}
}
