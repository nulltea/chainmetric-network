package rpc

import (
	"context"

	"github.com/timoth-y/chainmetric-network/orgservices/identity/api/middleware"
	"github.com/timoth-y/chainmetric-network/orgservices/identity/api/presenter/admin"
	presenter "github.com/timoth-y/chainmetric-network/orgservices/identity/api/presenter/user"
	"github.com/timoth-y/chainmetric-network/orgservices/identity/infrastructure/repository"
	"github.com/timoth-y/chainmetric-network/orgservices/identity/usecase/identity"
	"github.com/timoth-y/chainmetric-network/orgservices/identity/usecase/privileges"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/core"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type adminService struct{
	UnimplementedAdminServiceServer
}

// RegisterAdminService registers IdentityServiceServer fir given gRPC `server` instance.
func RegisterAdminService(server *grpc.Server) {
	RegisterAdminServiceServer(server, &adminService{})
}

func (s adminService) GetCandidates(ctx context.Context, _ *presenter.UsersRequest) (*presenter.UsersResponse, error) {
	var user = middleware.MustRetrieveUser(ctx)

	if !privileges.Has(user, "identity.enroll") {
		return nil, status.Error(codes.Unauthenticated, "user has not privileges for this method")
	}

	users, err := repository.NewUserMongo(core.MongoDB).ListByQuery(map[string]interface{}{
		"confirmed": false,
	})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "All users are confirmed already")
		}

		return nil, status.Errorf(codes.Internal, "failed to get users from database")
	}

	return presenter.NewUsersResponse(users), nil
}

// EnrollUser implements IdentityServiceClient gRPC service.
func (adminService) EnrollUser(
	ctx context.Context,
	request *admin.EnrollUserRequest,
) (*admin.EnrollUserResponse, error) {
	var user = middleware.MustRetrieveUser(ctx)

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

	return admin.NewEnrollUserResponse(initPassword), nil
}
