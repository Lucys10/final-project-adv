package handlers

import (
	"context"
	gouuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"hashserver/pkg/hashservice"
	"hashserver/pkg/logger"
	"hashserver/pkg/sha3hash"
	"sync"
)

type HashServer struct {
	Logs *logger.Log
	hashservice.UnimplementedHashServer
}

func (s HashServer) CalculateHash(ctx context.Context, listStr *hashservice.StrList) (*hashservice.HashList, error) {
	listStrToHash := listStr.Str
	var wg sync.WaitGroup
	listHash := make([]string, 0, len(listStrToHash))
	resCh := make(chan string)

	md, _ := metadata.FromIncomingContext(ctx)

	reqID := ""
	for _, v := range md.Get("ID") {
		if v == "" {
			uuid := gouuid.Must(gouuid.NewV4(), nil)
			reqID = uuid.String()
		}
	}

	for _, s := range listStrToHash {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			resCh <- sha3hash.GetHashSHA3(s)
		}(s)
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()

	for h := range resCh {
		listHash = append(listHash, h)
	}

	s.Logs.WithFields(logrus.Fields{
		"package":    "handlers",
		"handler":    "CalculateHash",
		"ID request": reqID,
		"quantity":   len(listHash),
	}).Info("Successful calculate hash")

	hash := &hashservice.HashList{Hash: listHash}
	return hash, nil
}
