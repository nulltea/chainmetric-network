package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/timoth-y/chainmetric-contracts/shared/core"
	"github.com/timoth-y/chainmetric-contracts/shared/infrastructure/repository"
	model "github.com/timoth-y/chainmetric-contracts/shared/model/user"
	"github.com/timoth-y/chainmetric-contracts/src/identity/usecase/access"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// JWTAuthUnaryInterceptor returns grpc.UnaryServerInterceptor func for JWT access control.
func JWTAuthUnaryInterceptor(skipMethods ...string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if skipValidation(info.FullMethod, skipMethods) {
			return handler(ctx, req)
		}

		user, err := tryRetrieveUserFromJWT(ctx)
		if err != nil {
			return nil, err
		}

		ctx = metadata.AppendToOutgoingContext(ctx,
			"user_id", user.ID,
			"user_model", user.Encode(),
		)

		return handler(ctx, req)
	}
}

// JWTAuthStreamInterceptor returns grpc.StreamServerInterceptor func for JWT access control.
func JWTAuthStreamInterceptor(skipMethods ...string) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if skipValidation(info.FullMethod, skipMethods) {
			return handler(srv, stream)
		}

		user, err := tryRetrieveUserFromJWT(stream.Context())
		if err != nil {
			return err
		}

		_ = stream.SetHeader(metadata.New(map[string]string{
			"user_id": user.ID,
			"user_model": user.Encode(),
		}))

		return handler(srv, stream)
	}
}

func tryRetrieveUserFromJWT(ctx context.Context) (*model.User, error) {
	meta, ok := metadata.FromIncomingContext(ctx); if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := meta["authorization"]; if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	token := values[0]
	claims, err := access.VerifyJWT(token); if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "access token is invalid: %w", err)
	}

	user, err := repository.NewUserMongo(core.MongoDB).GetByID(claims.Id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.Unauthenticated, "access token points to unknown user")
		}

		return nil, status.Errorf(codes.Internal, "failed to resolve user from token claims")
	}

	if !user.Confirmed {
		return nil, status.Errorf(codes.Unauthenticated, "user account hasn't been confirmed yet")
	}

	if user.ExpiresAt != nil && user.ExpiresAt.Before(time.Now()) {
		return nil, status.Errorf(codes.Unauthenticated, "user account is expired")
	}

	return user, nil
}

func skipValidation(method string, skipMethods []string) bool {
	skip := false

	if skipMethods == nil {
		return false
	}

	for i := range skipMethods {
		if fmt.Sprintf("/chainmetric.identity.service.%s", skipMethods[i]) == method {
			skip = true
			break
		}
	}

	return skip
}
