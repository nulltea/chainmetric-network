package rpc

import (
	"context"

	"github.com/timoth-y/chainmetric-contracts/shared/core"
	"github.com/timoth-y/chainmetric-contracts/shared/infrastructure/repository"
	"github.com/timoth-y/chainmetric-contracts/src/identity/api/middleware"
	"github.com/timoth-y/chainmetric-contracts/src/identity/api/presenter"
	"github.com/timoth-y/chainmetric-contracts/src/identity/usecase/access"
	"github.com/timoth-y/chainmetric-contracts/src/identity/usecase/identity"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userService struct{
	UnimplementedUserServiceServer
}

// RegisterUserService registers UserServiceServer fir given gRPC `server` instance.
func RegisterUserService(server *grpc.Server) {
	RegisterUserServiceServer(server, &userService{})
}

// Register implements IdentityServiceServer gRPC service.
func (userService) Register(
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

func (s userService) GetState(ctx context.Context, _ *emptypb.Empty) (*presenter.User, error) {
	var user = middleware.MustRetrieveUser(ctx)
	return presenter.NewUserProto(user), nil
}

// ChangePassword implements UserServiceServer gRPC service RPC.
func (userService) ChangePassword(
	ctx context.Context,
	request *presenter.ChangePasswordRequest,
) (*presenter.StatusResponse, error) {
	var user = middleware.MustRetrieveUser(ctx)

	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if user.Passcode != request.PrevPasscode {
		return nil, status.Error(codes.InvalidArgument, "previous passcode does not match")
	}

	if err := repository.NewUserMongo(core.MongoDB).UpdateByID(user.ID, map[string]interface{}{
		"passcode": request.NewPasscode,
	}); err != nil {
		return nil, status.Error(codes.Internal, "failed to update user in database")
	}

	if err := repository.NewIdentitiesVault(core.Vault).
		UpdateUserpassAccess(user.IdentityName(), request.NewPasscode); err != nil {
		return nil, status.Error(codes.Internal, "failed to update user pass on Vault")
	}

	return presenter.NewStatusResponse(presenter.Status_OK), nil
}
