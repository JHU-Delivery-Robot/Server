package main

import (
	"context"
	"fmt"
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
	"github.com/sirupsen/logrus"
)

func main() {
	config, err := LoadConfig("/etc/navserver/config.json")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	logger.Info("loading mTLS credentials...")

	credentials, err := middleware.LoadCredentials(
		config.Credentials.RootCA,
		config.Credentials.Certificate,
		config.Credentials.Key,
	)
	if err != nil {
		logger.Fatal(fmt.Errorf("credentials error: %v", err))
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		stopChannel := make(chan os.Signal, 1)
		signal.Notify(stopChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-stopChannel // blocking wait for signal
		logger.Info("received shutdown request, exiting...")
		cancel()
	}()

	storeLogger := logrus.WithField("service", "store")
	store, err := store.New(storeLogger)
	if err != nil {
		logger.Fatal(fmt.Errorf("store creation error: %v", err))
	}

	osrmClient := osrm.New(config.OSRMAddress, config.OSRMPRofileName)
	assigner := assigner.New(store, osrmClient)

	grpcLogger := logrus.WithField("service", "gRPC")
	grpcServer := grpcserver.New(store, assigner, credentials, config.GRPCListen, ctx, grpcLogger)
	restLogger := logrus.WithField("service", "REST")
	restServer := restserver.New(config.RESTListen, store, restLogger)

	errs := make(chan error)
	go func() {
		errs <- grpcServer.Run()
	}()
	go func() {
		errs <- restServer.Run(ctx)
	}()

	err = <-errs
	cancel()
	if err != nil {
		logger.Fatal(err)
	}
}
