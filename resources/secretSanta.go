package resources

import (
	// "errors"
	"github.com/gin-gonic/gin"
	// "github.com/iToto/jollyHelper/common"
	"github.com/iToto/jollyHelper/common/messagecode"
	"github.com/iToto/jollyHelper/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	// "strconv"
	"time"
)

type SecretSantaResource struct {
	Name       string
	Owner      string
	assignedOn int64
}

func (ss *SecretSantaResource) AssignNames(c *gin.Context) {
	var err error
	mongoStore := c.MustGet("mongoStore").(*mgo.Database)

	log.Printf("Ho Ho Ho, it's name choosing time!")

	// Get all names from db sorted by Age, Decreasing
	personModel := &models.Person{}
	persons := []models.Person{}
	personCollection := mongoStore.C(personModel.Collection())

	err = personCollection.EnsureIndex(personModel.Index())

	if err != nil {
		sendError(&err, messagecode.E_SERVER_ERROR, c)
		return
	}

	err = personCollection.Find(bson.M{}).Sort("-age").Limit(personModel.Limit()).All(&persons)
	// log.Printf("List of people: %s", persons)

	// Create names array list with
	exchangeList := make([]*SecretSantaResource, 0)
	for index, person := range persons {
		// Create new secret santa name for each person
		secretSantaName := &SecretSantaResource{}
		secretSantaName.Name = person.Name
		secretSantaName.Owner = ""
		secretSantaName.assignedOn = time.Now().Unix()
		exchangeList = append(exchangeList, secretSantaName)
	}

	log.Printf("Secret Santa List: %s", exchangeList)

	// TODO: Assign a name to each person
}
