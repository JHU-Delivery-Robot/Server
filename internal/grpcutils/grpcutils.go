package grpcutils

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

// SetupShutdown ensures clean server shutdown on system interrupts
func SetupShutdown(cancel context.CancelFunc, server *grpc.Server) {
	go func() {
		stopChannel := make(chan os.Signal, 1)
		signal.Notify(stopChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-stopChannel // blocking wait for signal
		cancel()
		server.GracefulStop()
	}()
}

// mergeServerContext creates new context which cancels when either gRPC client or server context cancels
func MergeServerContext(server_ctx context.Context, grpc_ctx context.Context) context.Context {
	new_ctx, cancel := context.WithCancel(grpc_ctx)

	go func() {
		select {
		case <-new_ctx.Done():
		case <-server_ctx.Done():
			cancel()
		}
	}()

	return new_ctx
}
