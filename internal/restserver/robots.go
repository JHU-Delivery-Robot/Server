package restserver

import (
	"encoding/json"
	"net/http"

	"github.com/JHU-Delivery-Robot/Server/internal/store"
	"github.com/gorilla/mux"
)

type Robot struct {
	ID       string        `json:"id"`
	Status   string        `json:"status"`
	Location store.Point   `json:"location"`
	Route    []store.Point `json:"route,omitempty"`
}

func (s *Server) getRobotsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		robots, err := s.store.GetRobots()
		if err != nil {
			s.logger.Error(err)
			Error(w, http.StatusInternalServerError)
			return
		}

		robotsAndRoutes := make([]Robot, 0)
		for _, robot := range robots {
			route, err := s.store.GetRoute(robot.ID)
			if err != nil {
				s.logger.Error(err)
				Error(w, http.StatusInternalServerError)
				return
			}

			robotsAndRoutes = append(robotsAndRoutes, Robot{
				ID:       robot.ID,
				Status:   robot.Status,
				Location: robot.Location,
				Route:    route,
			})
		}

		Respond(w, &robotsAndRoutes, http.StatusOK)
	}
}

func (s *Server) getRequestsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requests, err := s.store.GetAllRequests()
		if err != nil {
			s.logger.Error(err)
			Error(w, http.StatusInternalServerError)
			return
		}

		Respond(w, &requests, http.StatusOK)
	}
}

func (s *Server) postRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestLocation store.Point
		if err := json.NewDecoder(r.Body).Decode(&requestLocation); err != nil {
			Error(w, http.StatusUnprocessableEntity)
			return
		}

		requestID, err := s.store.CreateRequest(requestLocation)
		if err != nil {
			s.logger.Error(err)
			Error(w, http.StatusInternalServerError)
		}

		Respond(w, requestID, http.StatusCreated)
	}
}

func (s *Server) deleteRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		deleted, err := s.store.DeleteRequest(id)
		if err != nil {
			s.logger.Error(err)
			Error(w, http.StatusInternalServerError)
			return
		}

		if !deleted {
			Error(w, http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}