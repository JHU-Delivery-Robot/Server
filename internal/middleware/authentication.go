package middleware

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Authentication struct {
	robotPublicKeys map[string]string
}

const Identity = "identity"

func (a *Authentication) GetUnaryMiddleware() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("could not get metadata from incoming context")
		}

		md = md.Copy()
		md.Append(Identity, "Fangorn")

		ctx = metadata.NewIncomingContext(ctx, md)

		return handler(ctx, req)
	}
}
