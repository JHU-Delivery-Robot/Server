package routing

import (
	"context"

	pb "github.com/JHU-Delivery-Robot/Server/protocol"
)

type Route struct {
	Waypoints []*pb.Point
}

type Router interface {
	Route(ctx context.Context, start *pb.Point, end *pb.Point) (*Route, error)
}
