package router

import (
	"log/slog"
	"net/http"
	"orbit-app/internal/config"
	"orbit-app/internal/endpoints"
	"orbit-app/internal/middleware"
	"orbit-app/internal/services"
)

type Options struct {
	SnippetService services.SnippetService
	Logger         *slog.Logger
	Config         config.Config
}

func New(opt Options) http.Handler {
	router := http.NewServeMux()
	addRoutes(router, opt)
	handler := middleware.Pipeline(
		router,
		middleware.Logger(opt.Logger),
		middleware.StripSlash(),
		middleware.Recovery(opt.Logger),
	)
	return handler
}

func addRoutes(router *http.ServeMux, opt Options) {
	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("healthy"))
	})

	router.HandleFunc(
		"PUT /",
		endpoints.CreateSnippet(opt.SnippetService, opt.Logger, opt.Config.Server.Host).Unwrap(),
	)
	router.HandleFunc(
		"GET /{id}",
		endpoints.GetSnippetByID(opt.SnippetService, opt.Logger).Unwrap(),
	)
}
