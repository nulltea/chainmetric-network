package rpc

import (
	"context"

	"github.com/timoth-y/chainmetric-contracts/src/users/api/presenter"
	"github.com/timoth-y/chainmetric-contracts/src/users/usecase/identity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type identityService struct{}

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

func (identityService) Enroll(ctx context.Context, req *presenter.EnrollmentRequest) error {
	if err := req.Validate(); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if err := identity.Enroll(req.UserID,
		identity.WithRole(req.Role),
	); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}
