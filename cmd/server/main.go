package main

import (
	"context"
	"log"
	"net"

	"github.com/JHU-Delivery-Robot/Server/internal/grpcutils"
	"github.com/JHU-Delivery-Robot/Server/internal/middleware"
	"github.com/JHU-Delivery-Robot/Server/internal/osrm"
	"github.com/JHU-Delivery-Robot/Server/internal/server"
	pb "github.com/JHU-Delivery-Robot/Server/protocols"
	"google.golang.org/grpc"
)

func main() {
	config, err := LoadConfig("/etc/navserver/config.json")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	listener, err := net.Listen("tcp", ":443")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	authentication := middleware.Authentication{}
	credentials, err := authentication.LoadCredentials(
		config.Credentials.RootCA,
		config.Credentials.Certificate,
		config.Credentials.Key,
	)
	if err != nil {
		log.Fatalf("credentials error: %v", err)
	}

	grpc := grpc.NewServer(grpc.Creds(credentials), grpc.UnaryInterceptor(authentication.GetUnaryMiddleware()))
	server := server.New(ctx, osrm.New())
	pb.RegisterRoutingServer(grpc, &server)
	pb.RegisterDevelopmentServer(grpc, &server)

	grpcutils.SetupShutdown(cancel, grpc)

	if err := grpc.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
