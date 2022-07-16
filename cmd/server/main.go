package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/JHU-Delivery-Robot/Server/internal/grpcutils"
	"github.com/JHU-Delivery-Robot/Server/internal/middleware"
	"github.com/JHU-Delivery-Robot/Server/internal/osrm"
	"github.com/JHU-Delivery-Robot/Server/internal/server"
	pb "github.com/JHU-Delivery-Robot/Server/protocols"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func loadCredentials() (credentials.TransportCredentials, error) {
	tlsCert, err := tls.LoadX509KeyPair("certs/deliverbot_server.crt", "certs/deliverbot_server.key")
	if err != nil {
		return nil, fmt.Errorf("failed to load client key/cert pair: %v", err)
	}

	ca_data, err := ioutil.ReadFile("certs/deliverbot_ca.crt")
	if err != nil {
		return nil, fmt.Errorf("failed to load CA cert: %v", err)
	}

	ca_pool := x509.NewCertPool()
	if !ca_pool.AppendCertsFromPEM(ca_data) {
		return nil, fmt.Errorf("failed to add CA cert to pool: %v", err)
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{tlsCert},
		ClientCAs:    ca_pool,
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
	}

	return credentials.NewTLS(tlsConfig), nil
}

func main() {
	listener, err := net.Listen("tcp", ":443")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	authentication := middleware.Authentication{}
	credentials, err := loadCredentials()
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
