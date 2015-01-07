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
	"code.google.com/p/go.crypto/pbkdf2"
	"crypto/sha256"
)

// PERSONS_COL The collection name
const (
	PERSONS_COL       = "persons"
	MAX_GET_DOCUMENTS = 100 // The maximum number of documents to fetch
)

type ListItem struct {
	Uid         string `bson:"uid,omitempty" json:"uid" binding:"-"`
	Title       string `bson:"title,omitempty" json:"title,omitempty" binding:"required"`
	Description string `bson:"description,omitempty" json:"description,omitempty" binding:"required"`
}

type Login struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// JSON, ditch TitleCase for snake_case.
type Person struct {
	Uid       string     `bson:"uid,omitempty" json:"uid" binding:"-"`
	CreatedAt int64      `bson:"created_at,omitempty" json:"created_at" binding:"-"`
	UpdatedAt int64      `bson:"updated_at,omitempty" json:"updated_at" binding:"-"`
	Name      string     `bson:"name,omitempty" json:"name,omitempty" binding:"required"`
	Email     string     `bson:"email,omitempty" json:"email" binding:"required"`
	Password  string     `bson:"password,omitempty" json:"password" binding:"required"`
	Age       string     `bson:"age,omitempty" json:"age" binding:"required"`
	List      []ListItem `bson:"list,omitempty" json:"list" binding:"-"`
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

// Hash password
func (p *Person) HashPassword(password, salt []byte) string {
	bytePassword := pbkdf2.Key(password, salt, 4096, sha256.Size, sha256.New)

	return string(bytePassword[:])
}

func (p *Person) PasswordSalt() []byte {
	return []byte("l77t54l7")
}

// Collection getter method for the collections constant
func (p *Person) Collection() string {
	return PERSONS_COL
}

func (p *Person) Limit() int {
	return MAX_GET_DOCUMENTS
}
