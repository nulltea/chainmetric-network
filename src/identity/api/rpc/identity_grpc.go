package rpc

import (
	"context"

	"github.com/timoth-y/chainmetric-contracts/src/identity/api/presenter"
	"github.com/timoth-y/chainmetric-contracts/src/identity/usecase/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type identityService struct{
	UnimplementedIdentityServiceServer
}

// WithIdentityService registers IdentityServiceServer fir given gRPC `server` instance.
func WithIdentityService(server *grpc.Server) {
	RegisterIdentityServiceServer(server, &identityService{})
}

// Register implements IdentityServiceServer gRPC service.
func (identityService) Register(
	_ context.Context,
	request *presenter.RegistrationRequest,
) (*presenter.User, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := identity.Register(
		identity.WithName(request.Firstname, request.Lastname),
		identity.WithEmail(request.Email),
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return presenter.NewUserProto(*user), nil
}

// Enroll implements IdentityServiceClient gRPC service.
func (identityService) Enroll(
	_ context.Context,
	request *presenter.EnrollmentRequest,
) (*emptypb.Empty, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := identity.Enroll(request.UserID,
		identity.WithRole(request.Role),
		identity.WithExpirationPb(request.ExpireAt),
	); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}
