package grpcserver

import (
	"context"
)

// MergeContext creates new context which cancels when either parent cancels
func MergeContext(parent1 context.Context, parent2 context.Context) context.Context {
	child_ctx, cancel := context.WithCancel(parent1)

	// cancel derived if second parent ends
	go func() {
		<-parent2.Done()
		cancel()
	}()

	return child_ctx
}
