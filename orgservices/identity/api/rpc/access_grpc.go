package rpc

import (
	"context"
	"time"


	"github.com/timoth-y/chainmetric-network/orgservices/identity/api/presenter"
	"github.com/timoth-y/chainmetric-network/orgservices/identity/infrastructure/repository"
	"github.com/timoth-y/chainmetric-network/orgservices/identity/usecase/access"
	"github.com/timoth-y/chainmetric-network/orgservices/identity/usecase/identity"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/core"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type accessService struct{}

// RegisterAccessService registers AccessServiceServer fir given gRPC `server` instance.
func RegisterAccessService(server *grpc.Server) {
	RegisterAccessServiceServer(server, &accessService{})
}

// RequestFabricCredentials implements AccessServiceServer gRPC service RPC.
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

// AuthWithSigningIdentity implements AccessServiceServer gRPC service RPC.
func (s accessService) AuthWithSigningIdentity(
	ctx context.Context,
	request *presenter.CertificateAuthRequest,
) (*presenter.CertificateAuthResponse, error) {
	user, err := identity.ParseSigningCredentials(request.Certificate, request.SigningKey)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "identity verification failed")
	}

	accessToken, err := access.GenerateJWT(user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return presenter.NewCertificateAuthResponse(user, accessToken), nil
}
