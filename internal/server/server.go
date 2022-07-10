package server

import (
	"context"

	"github.com/JHU-Delivery-Robot/Server/internal/osrm"
	pb "github.com/JHU-Delivery-Robot/Server/protocols"
)

type routeSet map[string][]*pb.Point

type Server struct {
	pb.UnimplementedRoutingServer
	pb.UnimplementedDevelopmentServer
	osrm           osrm.Client
	routeOverrides routeSet
	ctx            context.Context
}

func New(ctx context.Context, osrm osrm.Client) Server {
	return Server{osrm: osrm, routeOverrides: make(routeSet), ctx: ctx}
}
