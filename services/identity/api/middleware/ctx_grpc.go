package middleware

import (
	"context"
	"fmt"

	model "github.com/timoth-y/chainmetric-network/shared/model/user"
	"google.golang.org/grpc/metadata"
)

// MustRetrieveUser returns 'user_model' from a given gRPC `ctx` if there is one, otherwise it panics.
func MustRetrieveUser(ctx context.Context) *model.User {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		panic("failed to get metadata from context")
	}

	if values := md.Get("user_model"); len(values) == 1 {
		return model.User{}.Decode(values[0])
	}

	panic("user_model is missing in context")
}

// TryRetrieveUser returns user model from a given gRPC `ctx` if there is one, otherwise returns error.
func TryRetrieveUser(ctx context.Context) (*model.User, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to get metadata from context")
	}

	if values := md.Get("user_model"); len(values) == 1 {
		return model.User{}.Decode(values[0]), nil
	}

	return nil, fmt.Errorf("user model is missing in context")
}

// MustRetrieveUserID returns `user_id` from a given gRPC `ctx` if there is one, otherwise it panics.
func MustRetrieveUserID(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		panic("failed to get metadata from context")
	}

	if values := md.Get("user_id"); len(values) == 1 {
		return values[0]
	}

	panic("user_id is missing in context")
}
