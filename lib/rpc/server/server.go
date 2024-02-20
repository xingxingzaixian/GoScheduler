package server

import (
	"GoScheduler/internal/modules/utils"
	pb "GoScheduler/lib/rpc/proto"
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"
)

type Server struct{}

var keepAlivePolicy = keepalive.EnforcementPolicy{
	MinTime:             10 * time.Second,
	PermitWithoutStream: true,
}

var keepAliveParams = keepalive.ServerParameters{
	MaxConnectionIdle: 30 * time.Second,
	Time:              30 * time.Second,
	Timeout:           3 * time.Second,
}

func (s Server) Run(ctx context.Context, req *pb.TaskRequest) (*pb.TaskResponse, error) {
	defer func() {
		if err := recover(); err != nil {
			zap.S().Error(err)
		}
	}()

	zap.S().Infof("execute cmd start: [id: %d cmd: %s]", req.Id, req.Command)
	output, err := utils.ExecShell(ctx, req.Command)

	resp := new(pb.TaskResponse)
	resp.Output = output
	if err != nil {
		resp.Error = err.Error()
	} else {
		resp.Error = ""
	}
	zap.S().Infof("execute cmd end: [id: %d cmd: %s err: %s]", req.Id, req.Command, resp.Error)

	return resp, nil
}

func Start(addr string) *grpc.Server {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		zap.S().Fatal(err)
	}

	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepAliveParams),
		grpc.KeepaliveEnforcementPolicy(keepAlivePolicy),
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterTaskServer(grpcServer, Server{})
	zap.S().Infof("server listen on %s", addr)

	go func() {
		err = grpcServer.Serve(l)
		if err != nil {
			zap.S().Fatal(err)
		}
	}()
	return grpcServer
}
