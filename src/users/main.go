package main

import (
	"net"

	"github.com/go-kit/kit/log/level"
	"github.com/timoth-y/chainmetric-contracts/shared/core"
	"github.com/timoth-y/chainmetric-contracts/shared/utils"
	"github.com/timoth-y/chainmetric-contracts/src/users/api/rpc"
	"github.com/timoth-y/chainmetric-contracts/src/users/usecase/identity"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
)

func init() {
	core.InitCore()
	core.InitMongoDB()
	utils.MustExecute(identity.Init, "failed to initialize identity package")
}

func main() {
	var (
		unaryInterceptors = []grpc.UnaryServerInterceptor{
			grpcTags.UnaryServerInterceptor(),
			grpcRecovery.UnaryServerInterceptor(),
		}

		streamInterceptors = []grpc.StreamServerInterceptor{
			grpcTags.StreamServerInterceptor(),
			grpcRecovery.StreamServerInterceptor(),
		}
	)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)

	rpc.IdentityServiceGRPC(server)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		return errors.Wrap(err, "failed to create listener")
	}

	_ = level.Info(core.Logger).Log("addr", s.grpcAddr, "msg", "Ready for gRPC")

	go func(lis net.Listener) {
		if err := s.gRPC.Serve(lis); err != nil {
			s.errs <- errors.Wrap(err, "failed to serve gRPC")
		}
	}(lis)
}
