package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type Authentication struct{}

const Identity = "identity"

func (a *Authentication) GetUnaryMiddleware() grpc.UnaryServerInterceptor {
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
