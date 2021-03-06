package resources

import (
	// "errors"
	"github.com/gin-gonic/gin"
	// "github.com/iToto/jollyHelper/common"
	"github.com/iToto/jollyHelper/common/messagecode"
	"github.com/iToto/jollyHelper/models"
	"github.com/mostafah/mandrill"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	// "strconv"
	// "encoding/json"
	// "reflect"
	"time"
)

type NotificationResource struct {
}

func (n *NotificationResource) Send(c *gin.Context) {
	// Get Secret Santa list

	id := c.Params.ByName("id")
	var err error
	mongoStore := c.MustGet("mongoStore").(*mgo.Database)

	log.Printf("Sending out secret Santa!")
	notification := &models.Notification{}
	notificationsCollection := mongoStore.C(notification.Collection())

	err = notificationsCollection.EnsureIndex(notification.Index())
	if err != nil {
		sendError(&err, messagecode.E_SERVER_ERROR, c)
		return
	}

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
	// log.Printf("List of people: %v", secretSantaList)

	// Create names array list with
	for _, secretSanta := range secretSantaList.List {
		// Send out secret santa for each person

		err = sendEmail(secretSanta.Owner.Email, secretSanta.Owner.Name, secretSanta.Name)

		if err != nil {
			sendError(&err, messagecode.E_SERVER_ERROR, c)
			return
		}

		// Save notification to DB
		notification.Uid = models.NewUid()
		notification.CreatedAt = time.Now().Unix()
		notification.Recipient = secretSanta.Owner.Email

		err = notificationsCollection.Insert(notification)
		if err != nil {
			sendError(&err, messagecode.E_SERVER_ERROR, c)
			return
		}

	}

}

func sendEmail(recipientEmail string, recipientName string, name string) error {
	// log.Printf("Sending email to: %v with name: %v", recipientEmail, name)
	msg := mandrill.NewMessageTo(recipientEmail, recipientName)
	msg.Text = "Ho Ho Ho, Merry Christmas " + recipientName + "! This year, your secret Santa is: " + name // optional
	msg.Subject = "2014 Secret Santa"
	msg.FromEmail = "noreply@jollyhelper.com"
	msg.FromName = "JollyHelper"
	res, err := msg.Send(false)

	log.Printf("Send email, response from Mandril: %v", res)

	if err != nil {
		return err
	}

	return nil
}
