package pprofroute

import (
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
)

func PprofRoute(r *chi.Mux) {
	r.Route("/debug/pprof", func(r chi.Router) {
		r.Get("/", pprof.Index)
		r.Get("/cpu", pprof.Profile)
		r.Get("/trace", pprof.Trace)
		r.Get("/heap", pprof.Handler("heap").ServeHTTP)
		r.Get("/goroutine", pprof.Handler("goroutine").ServeHTTP)
		r.Get("/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
		r.Get("/block", pprof.Handler("block").ServeHTTP)
		r.Get("/mutex", pprof.Handler("mutex").ServeHTTP)
	})
}
