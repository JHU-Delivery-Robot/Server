package assigner

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/JHU-Delivery-Robot/Server/internal/osrm"
	"github.com/JHU-Delivery-Robot/Server/internal/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type routeSet map[string][]store.Point

type Assigner struct {
	mu             sync.Mutex
	store          *store.Store
	routeOverrides routeSet
	osrm           osrm.Client
}

func New(store *store.Store, osrm osrm.Client) *Assigner {
	return &Assigner{
		store:          store,
		osrm:           osrm,
		routeOverrides: make(routeSet),
	}
}

func (a *Assigner) AddOverride(robotID string, route []store.Point) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.routeOverrides[robotID] = route
}

func (a *Assigner) Route(robotID string, robotLocation store.Point, ctx context.Context) ([]store.Point, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	route, err := a.store.GetRoute(robotID)
	if err != nil {
		return nil, fmt.Errorf("assigned route check: %v", err)
	}

	if route != nil {
		return route, err
	}

	waypoints, hasOverride := a.routeOverrides[robotID]
	if hasOverride {
		if err := a.store.AddRoute(robotID, waypoints, store.RequestNone); err != nil {
			return nil, err
		}
		return waypoints, nil
	}

	requests, err := a.store.GetIncompleteRequests()
	if err != nil {
		return nil, fmt.Errorf("getting incomplete requsts: %v", err)
	}

	for _, r := range requests {
		inProgress, err := a.store.IsInProgress(r.ID)
		if err != nil {
			return nil, fmt.Errorf("checking requests in-progress: %v", err)
		}

		if !inProgress {
			waypoints, err := a.osrm.Route(ctx, robotLocation, r.Location)
			if err != nil {
				log.Printf("error getting route: %v\n", err)
				return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get route: %v", err))
			}
			if err := a.store.AddRoute(robotID, waypoints, store.RequestNone); err != nil {
				return nil, err
			}
			return waypoints, nil
		}
	}

	return nil, nil
}
