package handlers

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"hashserver/pkg/hashservice"
	"hashserver/pkg/logger"
	"log"
	"net"
	"testing"
)

const bufsize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufsize)
	logs := logger.NewLogger(logrus.FatalLevel)
	s := grpc.NewServer()
	srv := &HashServer{Logs: logs}
	hashservice.RegisterHashServer(s, srv)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed start test GRPC Server: %v", err)
		}
	}()

}

func bufDialer(ctx context.Context, addr string) (net.Conn, error) {
	return lis.Dial()
}

func TestHashServer_CalculateHash(t *testing.T) {
	tests := map[string]struct {
		str  string
		hash string
	}{
		"One":   {str: "Mod1", hash: "d27be00225562b5d1130cc3e034009fd575cd48fe819b9114fa335ec98da9378"},
		"Two":   {str: "Mod2", hash: "14e84db48715d9d9197431c27de1c4836d0ebc9575db32f76fc77b746bf9388b"},
		"Three": {str: "Mod3", hash: "560baec0e5e885572b251394da0cbe9598d72819d0ab6c4fbc0bdd01977950ea"},
	}

	req := require.New(t)

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed connection GRPC dial: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("failed clone connection GRPC dial: %v", err)
		}
	}()

	client := hashservice.NewHashClient(conn)

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := client.CalculateHash(ctx, &hashservice.StrList{Str: []string{testCase.str}})
			if err != nil {
				log.Fatalf("failed CalculateHash")
			}

			for _, h := range resp.GetHash() {
				req.Equal(h, testCase.hash)
			}
		})

	}

}
