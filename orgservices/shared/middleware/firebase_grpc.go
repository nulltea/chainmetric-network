package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// FirebaseForUnaryGRPC returns grpc.UnaryServerInterceptor func for passing firebase registration token.
func FirebaseForUnaryGRPC(breakOnAbsence bool) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		token, ok := tryRetrieveFirebaseToken(ctx)

		if ok {
			ctx = metadata.AppendToOutgoingContext(ctx, "firebase_token", token)
		} else if breakOnAbsence {
			return nil, status.Error(codes.Unauthenticated, "firebase registration token missing")
		}

		return handler(ctx, req)
	}
}

// FirebaseForStreamGRPC returns grpc.StreamServerInterceptor func for passing firebase registration token.
func FirebaseForStreamGRPC(breakOnAbsence bool) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		token, ok := tryRetrieveFirebaseToken(stream.Context())

		if ok {
			_ = stream.SetHeader(metadata.New(map[string]string{"firebase_token": token}))
		} else if breakOnAbsence {
			return status.Error(codes.Unauthenticated, "firebase registration token missing")
		}

		return handler(srv, stream)
	}
}

func tryRetrieveFirebaseToken(ctx context.Context) (string, bool) {
	meta, ok := metadata.FromIncomingContext(ctx); if !ok {
		return "", false
	}

	values := meta["firebase_token"]; if len(values) == 0 {
		return "", false
	}

	return values[0], true
}
