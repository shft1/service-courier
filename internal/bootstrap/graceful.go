package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func StartServerGraceful(r chi.Router, pool *pgxpool.Pool, env *Env) {
	defer pool.Close()
	srv := &http.Server{Addr: "0.0.0.0:" + env.Port, Handler: r}

	srv.RegisterOnShutdown(func() {
		fmt.Println("Shutting down service-courier")
	})

	sysCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go srv.ListenAndServe()

	<-sysCtx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Println(err)
	}
}
