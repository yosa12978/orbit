package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"orbit-app/internal/config"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	conf   config.Config
	logger *slog.Logger
}

func New() *App {
	return &App{
		conf:   config.Get(),
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(
		context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
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
