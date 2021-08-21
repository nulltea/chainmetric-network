package rpc

import (
	"context"

	"github.com/timoth-y/chainmetric-contracts/shared/core"
	"github.com/timoth-y/chainmetric-contracts/shared/infrastructure/repository"
	"github.com/timoth-y/chainmetric-contracts/src/identity/api/presenter"
	"github.com/timoth-y/chainmetric-contracts/src/identity/usecase/auth"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type authService struct{
	UnimplementedAuthServiceServer
}

// RegisterAuthService registers AuthServiceServer fir given gRPC `server` instance.
func RegisterAuthService(server *grpc.Server) {
	RegisterAuthServiceServer(server, &authService{})
}

// Authenticate implements AuthServiceServer gRPC service RPC.
func (a authService) Authenticate(
	ctx context.Context,
	request *presenter.AuthRequest,
) (*presenter.AuthResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := repository.NewUserMongo(core.MongoDB).GetByQuery(map[string]interface{}{
		"email": request.Email,
		"password_hash": request.PasswordHash,
	})

	if err == mongo.ErrNoDocuments {
		return nil, status.Error(codes.InvalidArgument, "no user found for given credentials")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "failed to find user in database")
	}

	secretToken, err := auth.RequestVaultSecret(user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	accessToken, err := auth.GenerateJWT(user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return presenter.NewAuthResponse(secretToken, accessToken), nil
}

// SetPassword implements AuthServiceServer gRPC service RPC.
func (a authService) SetPassword(ctx context.Context, request *presenter.SetPasswordRequest) (*emptypb.Empty, error) {
	var userID = presenter.MustRetrieveUserID(ctx)

	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := repository.NewUserMongo(core.MongoDB).UpdateByID(userID, map[string]interface{}{
		"password_hash": request.PasswordHash,
	}); err != nil {
		return nil, status.Error(codes.Internal, "failed to update user in database")
	}

	return nil, nil
}
