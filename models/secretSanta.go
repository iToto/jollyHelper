package models

import (
	"gopkg.in/mgo.v2"
)

type SecretSanta struct {
	Uid        string `bson:"uid,omitempty" json:"uid" binding:"-"`
	Name       string `bson:"name,omitempty" json:"name,omitempty" binding:"required"`
	Owner      string `bson:"owner,omitempty" json:"owner" binding:"-"`
	AssignedOn int64  `bson:"assigned_on,omitempty" json:"assigned_on" binding:"-"`
}

const (
	SECRETSANTA_COL = "secret_santas"
)

// Index sets indexes on appropriate columns/properties
func (ss *SecretSanta) Index() mgo.Index {
	return mgo.Index{
		Key:    []string{"uid"},
		Unique: true,
		//DropDups:   true,
		Background: true,
		//Sparse:     true,
	}
}

// Collection getter method for the collections constant
func (ss *SecretSanta) Collection() string {
	return SECRETSANTA_COL
}
