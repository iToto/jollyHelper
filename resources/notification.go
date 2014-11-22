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
	// "encoding/json"
	// "reflect"
	// "time"
)

type NotificationResource struct {
}

func (n *NotificationResource) Send(c *gin.Context) {
	// Get Secret Santa list
	id := c.Params.ByName("id")
	var err error
	mongoStore := c.MustGet("mongoStore").(*mgo.Database)

	log.Printf("Sending out secret Santa!")

	// Get all names from db sorted by Age, Decreasing
	secretSantaModel := &models.SecretSanta{}
	secretSantaList := models.SecretSanta{}
	secretSantaCollection := mongoStore.C(secretSantaModel.Collection())

	err = secretSantaCollection.EnsureIndex(secretSantaModel.Index())

	if err != nil {
		sendError(&err, messagecode.E_SERVER_ERROR, c)
		return
	}

	err = secretSantaCollection.Find(bson.M{"uid": id}).One(&secretSantaList)
	// log.Printf("List of people: %s", secretSantaList)

	// Create names array list with
	for _, secretSanta := range secretSantaList.List {
		// TODO: Send out secret santa for each person

		log.Printf("Sending email to: %s with name: %s", secretSanta.Owner.Email, secretSanta.Name)
	}

}
