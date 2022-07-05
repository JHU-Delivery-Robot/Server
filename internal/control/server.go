package control

import (
	"context"

	"github.com/JHU-Delivery-Robot/Server/internal/osrm"
	pb "github.com/JHU-Delivery-Robot/Server/protocol"
)

type Server struct {
	pb.UnimplementedRoutingServer
}

func (s *Server) GetRoute(ctx context.Context, robotStatus *pb.RobotStatus) (*pb.Route, error) {
	destination := pb.Point{
		Latitude:  39.327864,
		Longitude: -76.6205217,
	}

	route, err := osrm.GetRoute(ctx, robotStatus.Location, &destination)
	return &pb.Route{Waypoints: route.Waypoints}, err
}
