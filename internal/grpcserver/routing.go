package grpcserver

import (
	"context"
	"fmt"

	"github.com/JHU-Delivery-Robot/Server/internal/middleware"
	"github.com/JHU-Delivery-Robot/Server/internal/store"
	pb "github.com/JHU-Delivery-Robot/Server/protocols"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ToProtocolRoute(route []store.Point) []*pb.Point {
	var waypoints = make([]*pb.Point, len(route))

	for i := 0; i < len(route); i++ {
		var waypoint pb.Point
		waypoint.Longitude = route[i].Longitude
		waypoint.Latitude = route[i].Latitude
		waypoints[i] = &waypoint
	}

	return waypoints
}

func (s *Server) GetRoute(clientCtx context.Context, robotStatus *pb.RobotStatus) (*pb.Route, error) {
	ctx := MergeContext(s.ctx, clientCtx)

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

	robot := store.Robot{
		ID:     identity[0],
		Status: robotStatus.Status.String(),
		Location: store.Point{
			Longitude: robotStatus.Location.Longitude,
			Latitude:  robotStatus.Location.Latitude,
		},
	}

	if err := s.store.UpsertRobot(robot); err != nil {
		s.logger.Error(fmt.Errorf("failed to save robot status: %v", err))
		return nil, status.Error(codes.Internal, "failed to save robot status")
	}

	waypoints, err := s.store.GetRoute(robot.ID)
	if err != nil {
		s.logger.Error(fmt.Errorf("could not assigned route: %v", err))
		return nil, status.Error(codes.Internal, "could not assigned route")
	}

	if waypoints != nil {
		pbWaypoints := ToProtocolRoute(waypoints)
		return &pb.Route{Waypoints: pbWaypoints}, nil
	}

	location := store.Point{
		Longitude: robotStatus.Location.Longitude,
		Latitude:  robotStatus.Location.Latitude,
	}
	waypoints, err = s.assigner.Route(robot.ID, location, ctx)
	if err != nil {
		s.logger.Error(fmt.Errorf("could not find route: %v", err))
		return nil, status.Error(codes.Internal, "could not find route")
	}

	pbWaypoints := ToProtocolRoute(waypoints)
	return &pb.Route{Waypoints: pbWaypoints}, nil
}
