package middleware

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const Identity = "identity"

func LoadCredentials(rootCA, serverCert, serverKey string) (credentials.TransportCredentials, error) {
	tlsCert, err := tls.LoadX509KeyPair(serverCert, serverKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load server key/cert pair: %v", err)
	}

	ca_data, err := ioutil.ReadFile(rootCA)
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

func MTLSHandler() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Internal, "could not get meta from incoming context")
		}

		peer, ok := peer.FromContext(ctx)
		if !ok {
			return nil, status.Error(codes.Internal, "could not get peer from incoming context")
		}

		mtls, ok := peer.AuthInfo.(credentials.TLSInfo)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "failed to get TLS auth info from peer")
		}

		if len(mtls.State.PeerCertificates) == 0 {
			return nil, status.Error(codes.Unauthenticated, "valid peer certificate not found")
		}

		client := mtls.State.PeerCertificates[0]
		md = md.Copy()
		md.Append(Identity, client.Subject.CommonName)
		ctx = metadata.NewIncomingContext(ctx, md)

		return handler(ctx, req)
	}
}
