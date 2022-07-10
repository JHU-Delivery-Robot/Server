package control

import (
	"context"

	"github.com/JHU-Delivery-Robot/Server/internal/grpcutils"
	"github.com/JHU-Delivery-Robot/Server/internal/routing"
	pb "github.com/JHU-Delivery-Robot/Server/protocol"
)

type Server struct {
	pb.UnimplementedRoutingServer
	router routing.Router
	ctx    context.Context
}

func NewServer(ctx context.Context, router routing.Router) Server {
	return Server{router: router, ctx: ctx}
}

func (s *Server) GetRoute(client_ctx context.Context, robotStatus *pb.RobotStatus) (*pb.Route, error) {
	ctx := grpcutils.MergeServerContext(s.ctx, client_ctx)

	destination := pb.Point{
		Latitude:  39.327864,
		Longitude: -76.6205217,
	}

	route, err := s.router.Route(ctx, robotStatus.Location, &destination)
	return &pb.Route{Waypoints: route.Waypoints}, err
}
