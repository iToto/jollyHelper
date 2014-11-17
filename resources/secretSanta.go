package resources

import (
	// "errors"
	"github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin/binding"
	// "github.com/iToto/jollyHelper/common"
	"github.com/iToto/jollyHelper/common/messagecode"
	"github.com/iToto/jollyHelper/models"
	"log"
	// "github.com/op/go-logging"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// "strconv"
	// "time"
)

type SecretSantaResource struct {
	Name       string
	Owner      string
	assignedOn int64
}

type Names struct {
}

func (ss *SecretSantaResource) AssignNames(c *gin.Context) {
	var err error
	mongoStore := c.MustGet("mongoStore").(*mgo.Database)

	log.Printf("Ho Ho Ho, it's name choosing time!")

	// Get all names from db sorted by Age, Decreasing
	person := &models.Person{}
	persons := []models.Person{}
	personCollection := mongoStore.C(person.Collection())

	err = personCollection.EnsureIndex(person.Index())

	if err != nil {
		sendError(&err, messagecode.E_SERVER_ERROR, c)
		return
	}

	err = personCollection.Find(bson.M{}).Sort("-age").Limit(person.Limit()).All(&persons)
	// log.Printf("List of people: %s", persons)

	// TODO: Create names array list with
	// TODO: Assign a name to each person
}
