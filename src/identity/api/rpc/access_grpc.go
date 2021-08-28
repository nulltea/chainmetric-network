package rpc

import (
	"context"
	"time"

	"github.com/timoth-y/chainmetric-contracts/shared/core"
	"github.com/timoth-y/chainmetric-contracts/shared/infrastructure/repository"
	"github.com/timoth-y/chainmetric-contracts/src/identity/api/presenter"
	"github.com/timoth-y/chainmetric-contracts/src/identity/usecase/access"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type accessService struct{
	UnimplementedAccessServiceServer
}

// RegisterAccessService registers AuthServiceServer fir given gRPC `server` instance.
func RegisterAccessService(server *grpc.Server) {
	RegisterAccessServiceServer(server, &accessService{})
}

// RequestFabricCredentials implements AuthServiceServer gRPC service RPC.
func (accessService) RequestFabricCredentials(
	_ context.Context,
	request *presenter.FabricCredentialsRequest,
) (*presenter.FabricCredentialsResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := repository.NewUserMongo(core.MongoDB).GetByQuery(map[string]interface{}{
		"email": request.Email,
		"passcode": request.Passcode,
	})

	if err == mongo.ErrNoDocuments {
		return nil, status.Error(codes.InvalidArgument, "no user found for given credentials")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "failed to find user in database")
	}

	if !user.Confirmed {
		return nil, status.Error(codes.Unavailable, "user in not confirmed by admin yet")
	}

	if user.ExpiresAt != nil && user.ExpiresAt.After(time.Now()) {
		return nil, status.Error(codes.Unavailable, "user account is expired")
	}

	secretPath, secretToken, err := access.RequestFabricCredentialsSecret(user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	accessToken, err := access.GenerateJWT(user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return presenter.NewFabricCredentialsResponse(user, secretToken, secretPath, accessToken), nil
}

// UpdatePassword implements AuthServiceServer gRPC service RPC.
func (accessService) UpdatePassword(
	ctx context.Context,
	request *presenter.UpdatePasswordRequest,
) (*presenter.StatusResponse, error) {
	var user = presenter.MustRetrieveUser(ctx)

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
