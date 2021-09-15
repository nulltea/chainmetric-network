package rpc

import (
	"context"

	"github.com/timoth-y/chainmetric-network/orgservices/notifications/api/presenter"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/proto"
	"google.golang.org/grpc"
)

type subscriberService struct{}

// RegisterSubscriberService registers SubscriberServiceServer fir given gRPC `server` instance.
func RegisterSubscriberService(server *grpc.Server) {
	RegisterSubscribeServiceServer(server, &subscriberService{})
}

func (subscriberService) Subscribe(
	ctx context.Context,
	request *presenter.SubscriptionRequest,
) (*proto.StatusResponse, error) {

	return nil, nil
}

func (subscriberService) Cancel(
	ctx context.Context,
	request *presenter.CancellationRequest,
) (*proto.StatusResponse, error) {

	return nil, nil
}
