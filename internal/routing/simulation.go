package routing

import (
	"context"

	pb "github.com/JHU-Delivery-Robot/Server/protocol"
)

type simulationRouter struct {
	waypoints []*pb.Point
}

func NewSimulationRouter(waypoints []*pb.Point) simulationRouter {
	return simulationRouter{waypoints: waypoints}
}

func (r *simulationRouter) Route(ctx context.Context, start *pb.Point, end *pb.Point) (*Route, error) {
	return &Route{Waypoints: r.waypoints}, nil
}
