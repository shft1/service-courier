package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/shft1/service-courier/observability/logger"
)

func StartPprofServer(ctx context.Context, log logger.Logger, r chi.Router) {
	srv := &http.Server{
		Addr:         net.JoinHostPort("localhost", "6060"),
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	srv.RegisterOnShutdown(func() {
		log.Info("shutting down pprof web-server...")
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
		log.Error("failed to start pprof web-server", logger.NewField("error", err))
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Warn("failed to shutdown pprof web-server gracefully", logger.NewField("error", err))
		return
	}
	log.Info("pprof web-server successfully stopped")
}
