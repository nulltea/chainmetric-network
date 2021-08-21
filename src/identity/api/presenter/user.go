package presenter

import (
	"context"

	model "github.com/timoth-y/chainmetric-contracts/shared/model/user"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewUserProto casts native user model to protobuf User.
func NewUserProto(user *model.User) *User {
	proto := &User{
		Id:        user.ID,
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

func MustRetrieveUser(ctx context.Context) *model.User {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		panic("failed to get metadata from context")
	}

	if values := md.Get("user_model"); len(values) == 1 {
		return model.User{}.Decode(values[0])
	}

	panic("user_model is missing in context")
}

func MustRetrieveUserID(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		panic("failed to get metadata from context")
	}

	if values := md.Get("user_id"); len(values) == 1 {
		return values[0]
	}

	panic("user_id is missing in context")
}
