package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/JHU-Delivery-Robot/Server/internal/control"
	"github.com/JHU-Delivery-Robot/Server/internal/grpcutils"
	"github.com/JHU-Delivery-Robot/Server/internal/routing"
	pb "github.com/JHU-Delivery-Robot/Server/protocol"
	"google.golang.org/grpc"
)

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [config path]\n", os.Args[0])
	}

	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	config, err := LoadConfig(args[0])
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	router := routing.NewSimulationRouter(config.GetRoute())

	server := grpc.NewServer()
	routing := control.NewServer(ctx, &router)
	pb.RegisterRoutingServer(server, &routing)

	grpcutils.SetupShutdown(cancel, server)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
