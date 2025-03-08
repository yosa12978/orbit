package app

import (
	"context"
	"io/fs"
	"log/slog"
	"net/http"
	"orbit-app/internal/config"
	"orbit-app/internal/data"
	"orbit-app/internal/repos"
	"orbit-app/internal/router"
	"orbit-app/internal/services"
)

func newServer(
	ctx context.Context,
	logger *slog.Logger,
	conf config.Config,
	assetsFS fs.FS,
) http.Server {
	snippetRepo := repos.NewSnippetRepoRedis(data.Redis(ctx), logger)
	snippetService := services.NewSnippetService(snippetRepo, logger)
	router := router.New(router.Options{
		SnippetService: snippetService,
		Logger:         logger,
		Config:         config.Get(),
		AssetsFS:       assetsFS,
	})
	return http.Server{
		Addr:           conf.Server.Addr,
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
	}
}
