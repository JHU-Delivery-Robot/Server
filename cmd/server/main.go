package main

import (
	"log"
	"net"

	"github.com/JHU-Delivery-Robot/Server/internal/control"
	pb "github.com/JHU-Delivery-Robot/Server/protocol"
	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRoutingServer(grpcServer, &control.Server{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
