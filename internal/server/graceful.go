package server

import (
	"context"
	"fmt"
	"net/http"
	"service-courier/internal/config/appcfg"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// StartServerGraceful - запуск сервера через graceful shutdown
func StartServerGraceful(ctx context.Context, r chi.Router, pool *pgxpool.Pool, env *appcfg.AppEnv) {
	srv := &http.Server{Addr: "0.0.0.0:" + env.AppPort, Handler: r}

	srv.RegisterOnShutdown(func() {
		fmt.Println("Shutting down service-courier")
	})

	go srv.ListenAndServe()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Println(err)
	}
}
