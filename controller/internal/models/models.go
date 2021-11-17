package models

type Hash struct {
	ID   int    `bson:"id" json:"id"`
	Hash string `bson:"hash" json:"hash"`
}
