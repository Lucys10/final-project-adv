package main

import (
	"context"
	"controller/internal/configs"
	"controller/internal/handlers"
	"controller/internal/store"
	"controller/pkg/db"
	"controller/pkg/hashservice"
	"controller/pkg/logger"
	"github.com/globalsign/mgo"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net/http"
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
		}).Fatal("failed to get configs")
	}

	connClient, err := grpc.DialContext(context.Background(), cfg.GrpcServerAddress,
		grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logs.WithFields(logrus.Fields{
			"package":  "main",
			"function": "grpcDialContext",
			"error":    err,
		}).Fatal("failed connection to grpc client")
	}

	hashClient := hashservice.NewHashClient(connClient)

	dbm, err := db.NewMongo(cfg.MongoURL)
	if err != nil {
		logs.WithFields(logrus.Fields{
			"package":  "main",
			"function": "NewMongo",
			"error":    err,
		}).Fatal("failed connection to mongodb")
	}
	s := store.NewStore(dbm)

	r := chi.NewRouter()
	h := handlers.Handlers{Ctx: context.Background(), GrpcClient: hashClient, Db: s, Logs: logs}
	h.RegisterRouter(r)

	server := &http.Server{Addr: cfg.ServerAddress, Handler: r}
	go func(s *http.Server) {
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logs.WithFields(logrus.Fields{
				"package":  "main",
				"function": "ListenAndServe",
				"error":    err,
			}).Fatal("The server is not up")
		}
	}(server)

	logs.WithFields(logrus.Fields{
		"ServerAddress": cfg.ServerAddress,
		"Log_level":     "Info",
	}).Info("Start controller-service...")

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	<-c

	if err := shutdown(server, connClient, dbm); err != nil {
		logs.WithFields(logrus.Fields{
			"package":  "main",
			"function": "shutdown",
			"error":    err,
		}).Fatal("failed shutdown service")
	}

}

func shutdown(server *http.Server, grpcClient *grpc.ClientConn, mongodb *mgo.Session) error {
	if err := server.Shutdown(context.Background()); err != nil {
		return errors.WithStack(err)
	}
	if err := grpcClient.Close(); err != nil {
		return errors.WithStack(err)
	}

	mongodb.Close()

	return nil
}
