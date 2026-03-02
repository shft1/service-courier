package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/shft1/service-courier/observability/logger"
)

// StartServer - запуск web-сервера через graceful shutdown
func StartServer(ctx context.Context, log logger.Logger, r chi.Router, host, port string) {
	srv := &http.Server{
		Addr:              net.JoinHostPort(host, port),
		Handler:           r,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       2 * time.Minute,
	}

	srv.RegisterOnShutdown(func() {
		log.Info("shutting down web-server...")
	})

	serverErr := make(chan error)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			serverErr <- err
			close(serverErr)
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
		log.Warn("failed to shutdown web-server gracefully", logger.NewField("error", err))
		return
	}
	log.Info("web-server successfully stopped")
}
