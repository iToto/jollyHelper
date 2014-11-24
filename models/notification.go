package models

import (
	// "github.com/iToto/jollyHelper/models"
	"gopkg.in/mgo.v2"
)

type Notification struct {
	Uid       string `bson:"uid,omitempty" json:"uid" binding:"-"`
	CreatedAt int64  `bson:"created_at,omitempty" json:"created_at" binding:"-"`
	Recipient string `bson:"recipient,omitempty" json:"recipient,omitempty" binding:"required"`
}

const (
	COLLECTION = "notifications"
)

// Index sets indexes on appropriate columns/properties
func (n *Notification) Index() mgo.Index {
	return mgo.Index{
		Key:    []string{"uid", "list"},
		Unique: true,
		//DropDups:   true,
		Background: true,
		//Sparse:     true,
	}
}

// Collection getter method for the collections constant
func (n *Notification) Collection() string {
	return COLLECTION
}
