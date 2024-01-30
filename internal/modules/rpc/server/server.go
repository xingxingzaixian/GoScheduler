package server

import (
	pb "GoScheduler/internal/modules/rpc/proto"
	"GoScheduler/internal/modules/utils"
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"os"
	"os/signal"
	"syscall"
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

func Start(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		zap.S().Fatal(err)
	}

	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepAliveParams),
		grpc.KeepaliveEnforcementPolicy(keepAlivePolicy),
	}
	server := grpc.NewServer(opts...)
	pb.RegisterTaskServer(server, Server{})
	zap.S().Infof("server listen on %s", addr)

	go func() {
		err = server.Serve(l)
		if err != nil {
			zap.S().Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	for {
		s := <-c
		zap.S().Infoln("收到信号 -- ", s)
		switch s {
		case syscall.SIGHUP:
			zap.S().Infoln("收到终端断开信号, 忽略")
		case syscall.SIGINT, syscall.SIGTERM:
			zap.S().Info("应用准备退出")
			server.GracefulStop()
			return
		}
	}
}
