package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"service-courier/observability/logger"
)

// StartServerGraceful - запуск сервера через graceful shutdown
func StartServerGraceful(ctx context.Context, log logger.Logger, r chi.Router, port string) {
	srv := &http.Server{
		Addr:              net.JoinHostPort("0.0.0.0:", port),
		Handler:           r,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       2 * time.Minute,
	}

	srv.RegisterOnShutdown(func() {
		log.Info("Shutting down service-courier...")
	})

	serverErr := make(chan error)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		log.Error("failed to start web-server", logger.NewField("error", err))
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Warn("failed to shutdown service-courier gracefully", logger.NewField("error", err))
		return
	}
	log.Info("service-courier successfully stopped")
}
