package store

import (
	"controller/internal/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
)

type Store interface {
	InsertHash(models.Hash) error
	GetHash(int) (models.Hash, error)
}

type store struct {
	db *mgo.Session
}

func NewStore(db *mgo.Session) *store {
	return &store{db: db}
}

func (s *store) InsertHash(hash models.Hash) error {
	if err := s.db.DB("hashdb").C("hash").Insert(hash); err != nil {
		return errors.WithStack(err)
	}
	return nil

}

func (s *store) GetHash(ids int) (models.Hash, error) {
	var hash models.Hash

	q := bson.M{
		"id": ids,
	}
	if err := s.db.DB("hashdb").C("hash").Find(q).One(&hash); err != nil {
		return models.Hash{}, errors.WithStack(err)
	}

	return hash, nil
}
