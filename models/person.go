package models

import (
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	// "code.google.com/p/go.crypto/bcrypt"
	// "code.google.com/p/go.crypto/pbkdf2"
	// "crypto/rand"
	// "crypto/sha256"
	// "github.com/iToto/jollyHelper/common"
	// "log"
)

// PERSONS_COL The collection name
const (
	PERSONS_COL       = "persons"
	MAX_GET_DOCUMENTS = 100 // The maximum number of documents to fetch
)

// JSON, ditch TitleCase for snake_case.
type Person struct {
	Uid       string `bson:"uid,omitempty" json:"uid" binding:"-"`
	CreatedAt int64  `bson:"created_at,omitempty" json:"created_at" binding:"-"`
	Name      string `bson:"name,omitempty" json:"name,omitempty" binding:"required"`
	Email     string `bson:"email,omitempty" json:"email" binding:"required"`
	Age       string `bson:"age,omitempty" json:"age" binding:"required"`
}

// Index sets indexes on appropriate columns/properties
func (p *Person) Index() mgo.Index {
	return mgo.Index{
		Key:    []string{"uid", "email"},
		Unique: true,
		//DropDups:   true,
		Background: true,
		//Sparse:     true,
	}
}

// Collection getter method for the collections constant
func (p *Person) Collection() string {
	return PERSONS_COL
}

func (p *Person) Limit() int {
	return MAX_GET_DOCUMENTS
}
