package main

import (
	"context"
	"log"
	"net"

	pb "github.com/JHU-Delivery-Robot/Server/spec"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedRouterServer
}

func (s *server) GetRoute(ctx context.Context, in *pb.Coords) (*pb.Route, error) {
	log.Printf("Received Coords: (%v, %v)\n", in.GetX(), in.GetY())
	return &pb.Route{Path: "No path available"}, nil
}

func main() {

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRouterServer(grpcServer, &server{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
