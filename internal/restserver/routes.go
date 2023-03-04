package restserver

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) registerRoutes() *mux.Router {
	r := mux.NewRouter()

	r.Handle("/docs/openapi.yaml", openAPIDocHandler()).Methods(http.MethodGet)

	r.Handle("/robots", s.getRobotsHandler()).Methods(http.MethodGet)
	r.Handle("/requests", s.getRequestsHandler()).Methods(http.MethodGet)
	r.Handle("/requests", s.postRequestHandler()).Methods(http.MethodPost)
	r.Handle("/requests/{id}", s.deleteRequestHandler()).Methods(http.MethodDelete)

	return r
}
