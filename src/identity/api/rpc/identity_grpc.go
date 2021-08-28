package rpc

import (
	"context"

	"github.com/timoth-y/chainmetric-contracts/shared/core"
	"github.com/timoth-y/chainmetric-contracts/shared/infrastructure/repository"
	"github.com/timoth-y/chainmetric-contracts/src/identity/api/presenter"
	"github.com/timoth-y/chainmetric-contracts/src/identity/usecase/access"
	"github.com/timoth-y/chainmetric-contracts/src/identity/usecase/identity"
	"github.com/timoth-y/chainmetric-contracts/src/identity/usecase/privileges"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type identityService struct{
	UnimplementedIdentityServiceServer
}

// RegisterIdentityService registers IdentityServiceServer fir given gRPC `server` instance.
func RegisterIdentityService(server *grpc.Server) {
	RegisterIdentityServiceServer(server, &identityService{})
}

// Register implements IdentityServiceServer gRPC service.
func (identityService) Register(
	_ context.Context,
	request *presenter.RegistrationRequest,
) (*presenter.RegistrationResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if _, err := repository.NewUserMongo(core.MongoDB).GetByQuery(map[string]interface{}{
		"email": request.Email,
	}); err != mongo.ErrNoDocuments {
		return nil, status.Errorf(codes.AlreadyExists,
			"user with email '%s' is already registered", request.Email)
	}

	user, err := identity.Register(
		identity.WithName(request.Firstname, request.Lastname),
		identity.WithEmail(request.Email),
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	accessToken, err := access.GenerateJWT(user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return presenter.NewRegistrationResponse(user, accessToken), nil
}

// Enroll implements IdentityServiceClient gRPC service.
func (identityService) Enroll(
	ctx context.Context,
	request *presenter.EnrollmentRequest,
) (*presenter.EnrollmentResponse, error) {
	var user = presenter.MustRetrieveUser(ctx)

	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if !privileges.Has(user, "identity.enroll") {
		return nil, status.Error(codes.Unauthenticated, "user has not privileges for this method")
	}

	initPassword, err := identity.Enroll(request.UserID,
		identity.WithRole(request.Role),
		identity.WithExpirationPb(request.ExpireAt),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return presenter.NewEnrollmentResponse(initPassword), nil
}
