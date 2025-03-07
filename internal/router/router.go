package router

import (
	"io/fs"
	"log/slog"
	"net/http"
	"orbit-app/internal/config"
	"orbit-app/internal/endpoints"
	"orbit-app/internal/middleware"
	"orbit-app/internal/services"
	"orbit-app/pkg"
)

type Options struct {
	SnippetService services.SnippetService
	Logger         *slog.Logger
	Config         config.Config
	AssetsFS       fs.FS
}

func New(opt Options) http.Handler {
	router := http.NewServeMux()
	addRoutes(router, opt)
	handler := middleware.Pipeline(
		router,
		middleware.Logger(opt.Logger),
		middleware.ContentType(opt.Logger),
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

	router.Handle(
		"GET /assets/",
		http.StripPrefix("/assets", http.FileServer(http.FS(opt.AssetsFS))),
	)

	router.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		pkg.RenderTemplate(w, "home", "Home page", nil)
	})

	router.HandleFunc(
		"PUT /",
		endpoints.CreateSnippet(
			opt.SnippetService,
			opt.Logger,
			opt.Config.Server.Host,
		).Unwrap(),
	)

	router.HandleFunc(
		"POST /",
		endpoints.CreateSnippet(
			opt.SnippetService,
			opt.Logger,
			opt.Config.Server.Host,
		).Unwrap(),
	)

	router.HandleFunc(
		"GET /{id}",
		endpoints.GetSnippetByID(
			opt.SnippetService,
			opt.Logger,
		).Unwrap(),
	)
}
