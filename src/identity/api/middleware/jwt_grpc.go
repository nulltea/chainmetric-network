package middleware

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func WithJWTAuthUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		metaOut, err := i.authenticate(ctx)
		if err != nil {
			return nil, err
		}

		ctx = metadata.NewOutgoingContext(ctx, metaOut)

		return handler(ctx, req)
	}
}

func WithJWTAuthStreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		metaOut, err := i.authenticate(stream.Context())
		if err != nil {
			return err
		}

		if err = stream.SetHeader(metaOut); err != nil {
			return status.Error(codes.Internal, "error setting metadata to outgoing context")
		}

		return handler(srv, stream)
	}
}

func tryRetrieveToken(ctx context.Context) (*meta.AuthToken, bool) {
	meta, ok := metadata.FromIncomingContext(ctx); if !ok {
		return "", nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := meta[AuthMetaKey]; if len(values) == 0 {
		return "", nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	claims, err := i.jwtManager.Verify(accessToken); if err != nil {
		return accessToken, nil, status.Errorf(codes.Unauthenticated, "access token is invalid: %w", err)
	}

	return accessToken, claims, nil
}
