package router

import (
	"log/slog"
	"net/http"
	"orbit-app/internal/middleware"
	"orbit-app/internal/services"
)

type Options struct {
	SnippetService services.SnippetService
	Logger         *slog.Logger
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
	router.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})

	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("healthy"))
	})
}
