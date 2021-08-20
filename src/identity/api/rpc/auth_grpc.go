package rpc

import (
	"context"

	"github.com/timoth-y/chainmetric-contracts/shared/core"
	"github.com/timoth-y/chainmetric-contracts/shared/infrastructure/repository"
	"github.com/timoth-y/chainmetric-contracts/src/identity/api/presenter"
	"github.com/timoth-y/chainmetric-contracts/src/identity/usecase/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type authService struct{
	UnimplementedAuthServiceServer
}

// WithAuthService registers AuthServiceServer fir given gRPC `server` instance.
func WithAuthService(server *grpc.Server) {
	RegisterAuthServiceServer(server, &authService{})
}

// Authenticate implements AuthServiceServer gRPC service.
func (a authService) Authenticate(
	ctx context.Context,
	request *presenter.AuthRequest,
) (*presenter.AuthResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := auth.Authenticate(request.Email, request.PasswordHash)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return presenter.NewAuthResponse(token), nil
}

func (a authService) SetPassword(ctx context.Context, request *presenter.SetPasswordRequest) (*emptypb.Empty, error) {
	var (
		repo = repository.NewUserMongo(core.MongoDB)
	)

	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := repo.GetByID(request.UserID)
	if err != nil {

	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return presenter.NewAuthResponse(token), nil
}
