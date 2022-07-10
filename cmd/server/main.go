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
	// flag.Usage = func() {
	// 	fmt.Printf("Usage: %s [config path]\n", os.Args[0])
	// }

	// flag.Parse()

	// args := flag.Args()
	// if len(args) != 1 {
	// 	flag.Usage()
	// 	os.Exit(1)
	// }

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	osrm := osrm.New()
	authentication := middleware.Authentication{}

	grpc := grpc.NewServer(grpc.UnaryInterceptor(authentication.GetUnaryMiddleware()))
	server := server.New(ctx, osrm)
	pb.RegisterRoutingServer(grpc, &server)
	pb.RegisterDevelopmentServer(grpc, &server)

	grpcutils.SetupShutdown(cancel, grpc)

	if err := grpc.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
