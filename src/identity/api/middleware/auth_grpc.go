package middleware

import (
	"context"
	"time"

	model "github.com/timoth-y/chainmetric-contracts/shared/model/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthForUnaryGRPC returns grpc.UnaryServerInterceptor func for JWT access control.
func AuthForUnaryGRPC(skipMethods ...string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if skipValidation(info.FullMethod, skipMethods) {
			return handler(ctx, req)
		}

		user, _ := TryRetrieveUser(ctx)

		// If previous interceptor skipped, skipping this one too:
		if user == nil {
			return handler(ctx, req)
		}

		if err := authenticateUser(user); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// AuthForStreamGRPC returns grpc.StreamServerInterceptor func for JWT access control.
func AuthForStreamGRPC(skipMethods ...string) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if skipValidation(info.FullMethod, skipMethods) {
			return handler(srv, stream)
		}

		user, _ := TryRetrieveUser(stream.Context())

		// If previous interceptor skipped, skipping this one too:
		if user == nil {
			return handler(srv, stream)
		}

		if err := authenticateUser(user); err != nil {
			return err
		}

		return handler(srv, stream)
	}
}

func authenticateUser(user *model.User) error {
	if !user.Confirmed {
		return status.Errorf(codes.Unauthenticated, "user account hasn't been confirmed yet")
	}

	if user.ExpiresAt != nil && user.ExpiresAt.Before(time.Now()) {
		return status.Errorf(codes.Unauthenticated, "user account is expired")
	}

	return nil
}
