package middleware

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/timoth-y/chainmetric-network/orgservices/identity/infrastructure/repository"
	"github.com/timoth-y/chainmetric-network/orgservices/identity/model"
	"github.com/timoth-y/chainmetric-network/orgservices/identity/usecase/access"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/core"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// JWTForUnaryGRPC returns grpc.UnaryServerInterceptor func for JWT access control.
func JWTForUnaryGRPC(skipMethods ...string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if skipValidation(info.FullMethod, skipMethods) {
			return handler(ctx, req)
		}

		claims, err := tryParseJWT(ctx)
		if err != nil {
			return nil, err
		}

		user, err := retrieveUserFromDB(claims.Id)
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

// JWTForStreamGRPC returns grpc.StreamServerInterceptor func for JWT access control.
func JWTForStreamGRPC(skipMethods ...string) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if skipValidation(info.FullMethod, skipMethods) {
			return handler(srv, stream)
		}

		claims, err := tryParseJWT(stream.Context())
		if err != nil {
			return err
		}

		user, err := retrieveUserFromDB(claims.Id)
		if err != nil {
			return err
		}

		_ = stream.SetHeader(metadata.New(map[string]string{
			"user_id":    claims.Id,
			"user_model": user.Encode(),
		}))

		return handler(srv, stream)
	}
}

func tryParseJWT(ctx context.Context) (*jwt.StandardClaims, error) {
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

	return claims, nil
}

func retrieveUserFromDB(id string) (*model.User, error) {
	user, err := repository.NewUserMongo(core.MongoDB).GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.Unauthenticated, "access token points to unknown user")
		}

		return nil, status.Errorf(codes.Internal, "failed to resolve user from token claims")
	}

	return user, nil
}

func skipValidation(method string, skipMethods []string) bool {
	skip := false

	if skipMethods == nil {
		return false
	}

	for i := range skipMethods {
		if fmt.Sprintf("/chainmetric.%s", skipMethods[i]) == method {
			skip = true
			break
		}
	}

	return skip
}
