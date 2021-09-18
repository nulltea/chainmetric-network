package rpc

import (
	"context"

	"github.com/timoth-y/chainmetric-network/orgservices/notifications/api/presenter"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/model/intention"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/usecase/eventproxy"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/middleware"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type subscriberService struct{}

// RegisterSubscriberService registers SubscriberServiceServer fir given gRPC `server` instance.
func RegisterSubscriberService(server *grpc.Server) {
	RegisterSubscriberServiceServer(server, &subscriberService{})
}

// Subscribe implements SubscriberServiceServer gRPC service RPC.
func (subscriberService) Subscribe(
	ctx context.Context,
	request *presenter.SubscriptionRequest,
) (_ *presenter.SubscriptionResponse, err error) {
	var (
		userToken = middleware.MustRetrieveFirebaseToken(ctx)
		concerns  []intention.EventConcern
	)

	if concerns, err = request.ToEventConcerns(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err = eventproxy.ForwardEvents(userToken, concerns...); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return presenter.NewSubscriptionResponse(concerns...), nil
}

// Cancel implements SubscriberServiceServer gRPC service RPC.
func (subscriberService) Cancel(
	ctx context.Context,
	request *presenter.CancellationRequest,
) (*proto.StatusResponse, error) {
	var userToken = middleware.MustRetrieveFirebaseToken(ctx)

	if err := eventproxy.CancelForwarding(userToken, request.Topics...); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}
