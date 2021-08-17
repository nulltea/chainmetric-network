package rpc

import (
	"context"

	"github.com/timoth-y/chainmetric-contracts/src/users/api/presenter"
	"github.com/timoth-y/chainmetric-contracts/src/users/usecase/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type identityService struct{
	UnimplementedIdentityServiceServer
}

func WithIdentityService(server *grpc.Server) {
	RegisterIdentityServiceServer(server, &identityService{})
}

// Register implements IdentityServiceClient gRPC service.
func (identityService) Register(
	ctx context.Context,
	req *presenter.RegistrationRequest,
) (*presenter.User, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := identity.Register(
		identity.WithName(req.Firstname, req.Lastname),
		identity.WithEmail(req.Email),
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return presenter.NewUserProto(*user), nil
}

// Enroll implements IdentityServiceClient gRPC service.
func (identityService) Enroll(ctx context.Context, req *presenter.EnrollmentRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := identity.Enroll(req.UserID,
		identity.WithRole(req.Role),
		identity.WithExpirationPb(req.ExpireAt),
	); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}
