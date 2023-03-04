package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/JHU-Delivery-Robot/Server/internal/assigner"
	"github.com/JHU-Delivery-Robot/Server/internal/grpcserver"
	"github.com/JHU-Delivery-Robot/Server/internal/middleware"
	"github.com/JHU-Delivery-Robot/Server/internal/osrm"
	"github.com/JHU-Delivery-Robot/Server/internal/restserver"
	"github.com/JHU-Delivery-Robot/Server/internal/store"
)

func main() {
	config, err := LoadConfig("/etc/navserver/config.json")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	credentials, err := middleware.LoadCredentials(
		config.Credentials.RootCA,
		config.Credentials.Certificate,
		config.Credentials.Key,
	)
	if err != nil {
		log.Fatalf("credentials error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		stopChannel := make(chan os.Signal, 1)
		signal.Notify(stopChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-stopChannel // blocking wait for signal
		cancel()
	}()

	store, err := store.New()
	if err != nil {
		log.Fatalf("store creation error: %v", err)
	}

	osrmClient := osrm.New(config.OSRMAddress, config.OSRMPRofileName)
	assigner := assigner.New(store, osrmClient)
	grpc_server := grpcserver.New(store, assigner, credentials, config.GRPCListen, ctx)
	rest_server := restserver.New(config.RESTListen, store)

	errs := make(chan error)
	go func() {
		errs <- grpc_server.Run()
	}()
	go func() {
		errs <- rest_server.Run(ctx)
	}()

	err = <-errs
	cancel()
	log.Fatalf("error: %v", err)
}
