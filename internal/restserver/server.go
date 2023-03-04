package restserver

import (
	"context"
	_ "embed"
	"log"
	"net/http"
	"time"

	"github.com/JHU-Delivery-Robot/Server/internal/store"
	"github.com/NYTimes/gziphandler"
)

// Server is the client-facing REST server
type Server struct {
	listenAddress string
	store         *store.Store
}

func New(listenAddress string, store *store.Store) Server {
	return Server{
		listenAddress: listenAddress,
		store:         store,
	}
}

func LimitBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 100e3) // 100 kB
		next.ServeHTTP(w, r)
	})
}

// Run starts the server
func (s *Server) Run(ctx context.Context) error {
	router := s.registerRoutes()

	var handler http.Handler = router
	handler = LimitBody(handler)
	handler = gziphandler.GzipHandler(handler)

	httpServer := &http.Server{
		Addr:              s.listenAddress,
		Handler:           handler,
		ReadTimeout:       time.Second * 15,
		ReadHeaderTimeout: time.Second * 15,
		WriteTimeout:      time.Second * 15,
		IdleTimeout:       time.Second * 30,
		MaxHeaderBytes:    4096,
	}

	errs := make(chan error)
	go func() {
		errs <- httpServer.ListenAndServe()
	}()

	log.Println("REST server listening...")

	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
		return httpServer.Shutdown(ctx)
	}
}

//go:embed openapi.yaml
var openAPI []byte

func openAPIDocHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		_, _ = w.Write(openAPI)
	}
}
