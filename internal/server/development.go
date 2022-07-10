package server

import (
	"context"
	"errors"
	"log"

	"github.com/JHU-Delivery-Robot/Server/internal/grpcutils"
	"github.com/JHU-Delivery-Robot/Server/internal/middleware"
	pb "github.com/JHU-Delivery-Robot/Server/protocols"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *Server) SetRoute(client_ctx context.Context, route *pb.Route) (*pb.RouteResponse, error) {
	ctx := grpcutils.MergeServerContext(s.ctx, client_ctx)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("could not get meta from incoming context")
	}

	identity := md.Get(middleware.Identity)
	if len(identity) != 1 {
		log.Println("no identity found")
		return nil, status.Error(codes.Unauthenticated, "no identity provided")
	}

	s.AddRoute(identity[0], route.Waypoints)
	log.Printf("route set for: %s\n", identity[0])

	return &pb.RouteResponse{}, nil
}
