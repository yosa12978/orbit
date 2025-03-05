package app

import (
	"context"
	"errors"
	"io/fs"
	"log/slog"
	"net/http"
	"orbit-app/internal/config"
	"orbit-app/pkg"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	conf      config.Config
	logger    *slog.Logger
	templates fs.FS
}

func New(templates fs.FS) *App {
	return &App{
		conf:      config.Get(),
		logger:    slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		templates: templates,
	}
}

func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(
		context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	pkg.InitTemplates(a.templates)
	server := newServer(ctx, a.logger, a.conf)
	errCh := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil &&
			errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()
	var err error
	select {
	case err = <-errCh:
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(
			context.Background(), 10*time.Second)
		defer cancel()
		err = server.Shutdown(timeout)
	}
	return err
}
