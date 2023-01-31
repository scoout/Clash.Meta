package route

import (
	"github.com/Dreamacro/clash/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"runtime"
)

func debugRouter() http.Handler {
	handler := middleware.Profiler()
	r := chi.NewRouter()
	r.Mount("/", handler)
	r.Put("/gc", func(writer http.ResponseWriter, request *http.Request) {
		log.Debugln("trigger GC")
		runtime.GC()
	})

	return r
}
