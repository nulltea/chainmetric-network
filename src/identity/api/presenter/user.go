package presenter

import (
	model "github.com/timoth-y/chainmetric-contracts/shared/model/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewUserProto casts native user model to protobuf User.
func NewUserProto(user *model.User) *User {
	proto := &User{
		Id:        user.ID,
		Username:  user.IdentityName(),
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Role:      user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		Confirmed: user.Confirmed,
	}

	if user.ExpiresAt != nil {
		proto.ExpireAt = timestamppb.New(*user.ExpiresAt)
	}

	return proto
}

// NewUsersResponse presents UsersResponse for given native models slice `users`.
func NewUsersResponse(users []*model.User) *UsersResponse {
	var resp = &UsersResponse{
		Count: int64(len(users)),
	}

	for i := range users {
		resp.Users = append(resp.Users, NewUserProto(users[i]))
	}

	return resp
}

// NewRegistrationResponse presents RegistrationResponse for gRPC proto for given `user`,
// and grants access via `jwt`.
func NewRegistrationResponse(user *model.User, jwt string) *RegistrationResponse {
	return &RegistrationResponse{
		User: NewUserProto(user),
		AccessToken: jwt,
	}
}

// NewUserStatusResponse presents UserStatusResponse for gRPC proto for given `user`,
// and grants access via `jwt`.
func NewUserStatusResponse(status model.Status, initialPassword *string) *UserStatusResponse {
	return &UserStatusResponse{
		Status: UserStatus(status),
		InitialPassword: initialPassword,
	}
}

