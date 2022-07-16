package server

import (
	"context"
	"fmt"
	"log"

	"github.com/JHU-Delivery-Robot/Server/internal/grpcutils"
	"github.com/JHU-Delivery-Robot/Server/internal/middleware"
	pb "github.com/JHU-Delivery-Robot/Server/protocols"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *Server) AddRoute(name string, waypoints []*pb.Point) {
	s.routeOverrides[name] = waypoints
}

func (s *Server) GetRoute(client_ctx context.Context, robotStatus *pb.RobotStatus) (*pb.Route, error) {
	ctx := grpcutils.MergeServerContext(s.ctx, client_ctx)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "could not get meta from incoming context")
	}

	identity := md.Get(middleware.Identity)
	if len(identity) != 1 {
		log.Printf("valid identity not found\n")
		return nil, status.Error(codes.Unauthenticated, "invalid authentication credentials")
	}

	waypoints, hasOverride := s.routeOverrides[identity[0]]
	if hasOverride {
		return &pb.Route{Waypoints: waypoints}, nil
	}

	destination := pb.Point{
		Latitude:  39.327864,
		Longitude: -76.6205217,
	}

	waypoints, err := s.osrm.Route(ctx, robotStatus.Location, &destination)
	if err != nil {
		log.Printf("error getting route: %v\n", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get route: %v", err))
	}

	return &pb.Route{Waypoints: waypoints}, err
}
