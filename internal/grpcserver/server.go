package grpcserver

import (
	"context"
	"fmt"
	"net"

	"github.com/JHU-Delivery-Robot/Server/internal/assigner"
	"github.com/JHU-Delivery-Robot/Server/internal/middleware"
	"github.com/JHU-Delivery-Robot/Server/internal/store"
	pb "github.com/JHU-Delivery-Robot/Server/protocols"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	pb.UnimplementedRoutingServer
	pb.UnimplementedDevelopmentServer

	mTLSCredentials credentials.TransportCredentials
	listenAddress   string

	logger   *logrus.Entry
	store    *store.Store
	assigner *assigner.Assigner

	ctx context.Context
}

func New(
	store *store.Store,
	assigner *assigner.Assigner,
	mTLSCredentials credentials.TransportCredentials,
	listenAddress string,
	ctx context.Context,
	logger *logrus.Entry,
) Server {
	return Server{
		mTLSCredentials: mTLSCredentials,
		listenAddress:   listenAddress,
		logger:          logger,
		store:           store,
		assigner:        assigner,
		ctx:             ctx,
	}
}

func (s *Server) Run() error {
	grpc := grpc.NewServer(grpc.Creds(s.mTLSCredentials), grpc.UnaryInterceptor(middleware.MTLSHandler()))
	pb.RegisterRoutingServer(grpc, s)
	pb.RegisterDevelopmentServer(grpc, s)

	listener, err := net.Listen("tcp", s.listenAddress)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	errs := make(chan error)
	go func() {
		errs <- grpc.Serve(listener)
	}()

	s.logger.Info("gRPC server listening...")

	select {
	case err := <-errs:
		return err
	case <-s.ctx.Done():
		grpc.GracefulStop()
		return nil
	}
}
