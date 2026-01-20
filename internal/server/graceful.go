package server

import (
	"context"
	"net/http"
	"service-courier/observability/logger"
	"time"

	"github.com/go-chi/chi/v5"
)

// StartServerGraceful - запуск сервера через graceful shutdown
func StartServerGraceful(ctx context.Context, log logger.Logger, r chi.Router, port string) {
	srv := &http.Server{Addr: "0.0.0.0:" + port, Handler: r}

	srv.RegisterOnShutdown(func() {
		log.Info("Shutting down service-courier...")
	})

	go srv.ListenAndServe()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("failed to shutdown service-courier", logger.NewField("error", err))
		return
	}
	log.Info("service-courier successfully stopped")
}
