package db

import (
	"github.com/globalsign/mgo"
	"github.com/pkg/errors"
)

func NewMongo(url string) (*mgo.Session, error) {
	c, err := mgo.Dial(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := c.Ping(); err != nil {
		return nil, errors.WithStack(err)
	}

	return c, nil
}
