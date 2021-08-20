package presenter

import (
	model "github.com/timoth-y/chainmetric-contracts/shared/model/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewUserProto casts native user model to protobuf User.
func NewUserProto(user model.User) *User {
	proto := &User{
		Id:        user.ID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Role:      user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		Confirmed: user.Confirmed,
	}

	if user.ExpireAt != nil {
		proto.ExpireAt = timestamppb.New(*user.ExpireAt)
	}

	return proto
}
