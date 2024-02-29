package main

import (
	"log/slog"
	"sync"

	"github.com/kiennyo/syncwatch-be/internal/domain/auth"
	"github.com/kiennyo/syncwatch-be/internal/infrastructure/config"
	"github.com/kiennyo/syncwatch-be/internal/infrastructure/http"
	"github.com/kiennyo/syncwatch-be/internal/infrastructure/log"
)

var wg sync.WaitGroup

func main() {
	log.Init()

	server := http.New(&wg, config.HTTP{
		Port: 3000,
	}).
		AddRoutes("/auth", auth.Handlers())

	if err := server.Serve(); err != nil {
		slog.Error("Failed to start server", "reason", err.Error()) // Fatal
	}
}
