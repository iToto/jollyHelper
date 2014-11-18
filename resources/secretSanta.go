package resources

import (
	"errors"
	"github.com/gin-gonic/gin"
	// "github.com/iToto/jollyHelper/common"
	"github.com/iToto/jollyHelper/common/messagecode"
	"github.com/iToto/jollyHelper/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	// "strconv"
	"math/rand"
	"time"
)

type SecretSantaResource struct {
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
	exchangeList := make([]models.SecretSanta, 0)
	for _, person := range persons {
		// Create new secret santa name for each person
		secretSantaName := models.SecretSanta{}
		secretSantaName.Name = person.Name
		secretSantaName.Owner = ""
		secretSantaName.AssignedOn = 0
		exchangeList = append(exchangeList, secretSantaName)
	}

	log.Printf("Secret Santa List: %s", exchangeList)

	// Assign a name to each person
	err = runSecretSanta(persons, exchangeList)

	if err != nil {
		sendError(&err, messagecode.E_SERVER_ERROR, c)
	}

	log.Printf("Secret Santa List Done: %s", exchangeList)
}

func runSecretSanta(persons []models.Person, secretSantaList []models.SecretSanta) error {
	// For every person, assign them a secret santa

	if len(persons) == 0 || len(secretSantaList) == 0 {
		return errors.New("Arrays cannot be empty")
	}

	rand.Seed(time.Now().UTC().UnixNano())
	for index, person := range persons {
		// Loop until unused name is chosen that is not current person
		var i int
		for i = randInt(0, len(persons)); (i == index) || (secretSantaList[i].AssignedOn != 0); i = randInt(0, len(persons)) {
			log.Printf("Re-choosing i: %i - AssignedOn: %i", i, secretSantaList[i].AssignedOn)
		}
		// Selected secret santa, update secretSanta list with secret santa
		log.Printf("Name selected for %s", person.Name)
		secretSantaList[i].AssignedOn = time.Now().Unix()
		secretSantaList[i].Owner = person.Name
	}
	return nil
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
