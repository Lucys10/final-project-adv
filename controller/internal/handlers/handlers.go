package handlers

import (
	"context"
	"controller/internal/models"
	"controller/internal/store"
	"controller/pkg/generateid"
	"controller/pkg/hashservice"
	"controller/pkg/logger"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"io"
	"net/http"
	"strconv"
)

type Handlers struct {
	Ctx        context.Context
	GrpcClient hashservice.HashClient
	Db         store.Store
	Logs       *logger.Log
}

func (h *Handlers) RegisterRouter(r *chi.Mux) {
	r.Post("/send", h.Send)
	r.Get("/check", h.Check)
}

func (h *Handlers) Send(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Logs.WithFields(logrus.Fields{
			"package":  "handlers",
			"handler":  "send",
			"function": "ReadAll(r.body)",
			"error":    err,
		}).Error("failed read Body HTTP request")
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			h.Logs.WithFields(logrus.Fields{
				"package":  "handlers",
				"handler":  "send",
				"function": "Body.Close",
				"error":    err,
			}).Error("failed close Body HTTP request")
		}
	}()

	var arrayOfStrings []string
	if err := json.Unmarshal(data, &arrayOfStrings); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Logs.WithFields(logrus.Fields{
			"package":  "handlers",
			"handler":  "send",
			"function": "Unmarshal arrayOfStrings",
			"error":    err,
		}).Error("failed unmarshal arrayOfStrings")
		return
	}

	listStrToHash := &hashservice.StrList{Str: arrayOfStrings}

	ctx := requestID("X-Request_ID", r)
	hashList, err := h.GrpcClient.CalculateHash(ctx, listStrToHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Logs.WithFields(logrus.Fields{
			"package":  "handlers",
			"handler":  "send",
			"function": "GrpcClient CalculateHash",
			"error":    err,
		}).Error("failed request to Grpc server")
		return
	}

	arrayOfHash := make([]models.Hash, 0, len(arrayOfStrings))
	for _, hs := range hashList.GetHash() {
		hash := models.Hash{}
		hash.ID = generateid.GenerateID()
		hash.Hash = hs
		arrayOfHash = append(arrayOfHash, hash)
		if err := h.Db.InsertHash(hash); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.Logs.WithFields(logrus.Fields{
				"package":  "handlers",
				"handler":  "send",
				"function": "Db InsertHash",
				"error":    err,
			}).Error("failed insert hash to db")
			return
		}
	}

	resp, err := json.Marshal(arrayOfHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Logs.WithFields(logrus.Fields{
			"package":  "handlers",
			"handler":  "send",
			"function": "Marshal",
			"error":    err,
		}).Error("failed marshal arrayOfHash")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Logs.WithFields(logrus.Fields{
			"package":  "handlers",
			"handler":  "send",
			"function": "Write response",
			"error":    err,
		}).Error("failed write response")
		return
	}
}

func (h *Handlers) Check(w http.ResponseWriter, r *http.Request) {
	ids, ok := r.URL.Query()["ids"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		h.Logs.WithFields(logrus.Fields{
			"package":  "handlers",
			"handler":  "check",
			"function": "Url query",
		}).Errorf("Don't query parametr")
		return
	}

	idsInt, err := strconv.Atoi(ids[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Logs.WithFields(logrus.Fields{
			"package":  "handlers",
			"handler":  "check",
			"function": "Atoi",
			"error":    err,
		}).Error("failed convert str to int")
		return
	}

	res, err := h.Db.GetHash(idsInt)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		h.Logs.WithFields(logrus.Fields{
			"package":  "handlers",
			"handler":  "check",
			"function": "Db GetHash",
			"error":    err,
		}).Error("failed get hash with db")
		return
	}

	resp, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Logs.WithFields(logrus.Fields{
			"package":  "handlers",
			"handler":  "check",
			"function": "Marshal",
			"error":    err,
		}).Error("failed marshal hash")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Logs.WithFields(logrus.Fields{
			"package":  "handlers",
			"handler":  "check",
			"function": "Write",
			"error":    err,
		}).Error("failed write response")
		return
	}
}

func requestID(key string, r *http.Request) context.Context {
	reqID := r.Header.Get(key)
	ctx := metadata.AppendToOutgoingContext(context.Background(), "ID", reqID)
	return ctx
}
