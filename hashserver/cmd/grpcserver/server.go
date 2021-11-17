package main

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"hashserver/internal/configs"
	"hashserver/internal/handlers"
	"hashserver/pkg/hashservice"
	"hashserver/pkg/logger"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logs := logger.NewLogger(logrus.InfoLevel)

	cfg, err := configs.GetConfig()
	if err != nil {
		logs.WithFields(logrus.Fields{
			"package":  "main",
			"function": "GetConfig",
			"error":    err,
		}).Fatal("failed get to configs")
	}

	lis, err := net.Listen(cfg.Network, cfg.GrpcServerAddress)
	if err != nil {
		logs.WithFields(logrus.Fields{
			"package":  "main",
			"function": "net Listen",
			"error":    err,
		}).Fatal("failed create network")
	}

	s := grpc.NewServer()
	srv := &handlers.HashServer{Logs: logs}
	hashservice.RegisterHashServer(s, srv)

	go func(s *grpc.Server, lis net.Listener) {
		if err := s.Serve(lis); err != nil {
			logs.WithFields(logrus.Fields{
				"package":  "main",
				"function": "Serve",
				"error":    err,
			}).Fatal("The server is not up")
		}
	}(s, lis)

	logs.WithFields(logrus.Fields{
		"GRPC Server Address": cfg.GrpcServerAddress,
		"Log_level":           "Info",
	}).Info("Start hash-server...")

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	<-c

	s.GracefulStop()
}
