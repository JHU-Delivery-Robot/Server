package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func (s *server) router(ctx context.Context, in *pb.Coords) (*pb.Route, error) {
	log.Printf("Received Coords: (%v, %v)\n", in.GetX(), in.GetY())
	return &pb.Route{path: "No path available"}, nil
}
