package models

import (
	// "github.com/iToto/jollyHelper/models"
	"gopkg.in/mgo.v2"
)

type NameEntry struct {
	Name       string
	Owner      Person
	AssignedOn int64
}

type SecretSanta struct {
	Uid       string      `bson:"uid,omitempty" json:"uid" binding:"-"`
	CreatedAt int64       `bson:"created_at,omitempty" json:"created_at" binding:"-"`
	List      []NameEntry `bson:"list,omitempty" json:"list,omitempty" binding:"required"`
}

const (
	SECRETSANTA_COL = "secret_santas"
)

// Index sets indexes on appropriate columns/properties
func (ss *SecretSanta) Index() mgo.Index {
	return mgo.Index{
		Key:    []string{"uid", "list"},
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
