package control

import (
	"context"
	"log"

	pb "github.com/JHU-Delivery-Robot/Server/spec"
)

type Server struct {
	pb.UnimplementedRouterServer
}

func (s *Server) GetRoute(ctx context.Context, in *pb.Coords) (*pb.Route, error) {
	log.Printf("Received Coords: (%v, %v)\n", in.GetX(), in.GetY())
	return &pb.Route{Path: "No path available"}, nil
}

func (s *Server) SendStatus(ctx context.Context, in *pb.RobotStatus) (*pb.StatusConfirmation, error) {
	log.Print("Recieved robot status\n")
	return &pb.StatusConfirmation{Errors: 0}, nil
}
