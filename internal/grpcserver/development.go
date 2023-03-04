package grpcserver

import (
	"context"

	"github.com/JHU-Delivery-Robot/Server/internal/middleware"
	"github.com/JHU-Delivery-Robot/Server/internal/store"
	pb "github.com/JHU-Delivery-Robot/Server/protocols"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ToStoreRoute(route []*pb.Point) []store.Point {
	var waypoints = make([]store.Point, len(route))

	for i := 0; i < len(route); i++ {
		waypoints[i] = store.Point{
			Longitude: route[i].Longitude,
			Latitude:  route[i].Latitude,
		}
	}

	return waypoints
}

func (s *Server) SetRoute(clientCtx context.Context, route *pb.Route) (*pb.RouteResponse, error) {
	md, ok := metadata.FromIncomingContext(clientCtx)
	if !ok {
		s.logger.Error("could not get meta from incoming context")
		return nil, status.Error(codes.Internal, "could not get meta from incoming context")
	}

	identity := md.Get(middleware.Identity)
	if len(identity) != 1 {
		s.logger.Error("valid identity not found")
		return nil, status.Error(codes.Unauthenticated, "invalid authentication credentials")
	}

	s.assigner.AddOverride(identity[0], ToStoreRoute(route.Waypoints))
	s.logger.Infof("route override set for: %s", identity[0])

	return &pb.RouteResponse{}, nil
}
