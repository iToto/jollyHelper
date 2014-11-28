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
	// "encoding/json"
	"math/rand"
	"reflect"
	"time"
)

type SecretSantaResource struct {
}

func (ss *SecretSantaResource) List(c *gin.Context) {
	var err error
	mongoStore := c.MustGet("mongoStore").(*mgo.Database)
	secretSantaModel := &models.SecretSanta{}
	secretSantaList := []models.SecretSanta{}
	secretSantaCollection := mongoStore.C(secretSantaModel.Collection())

	err = secretSantaCollection.EnsureIndex(secretSantaModel.Index())
	if err != nil {
		sendError(&err, messagecode.E_SERVER_ERROR, c)
		return
	}

	err = secretSantaCollection.Find(bson.M{}).Sort("-created_at").All(&secretSantaList)
	if err != nil {
		if err == mgo.ErrNotFound {
			sendError(&err, messagecode.S_RESOURCE_NOTFOUND, c)
		} else {
			sendError(&err, messagecode.E_SERVER_ERROR, c)
		}
		return
	}

	for index, _ := range secretSantaList {
		secretSantaList[index].List = make([]models.NameEntry, 0)
	}
	sendResponse(secretSantaList, messagecode.S_RESOURCE_OK, c)
	return
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

	// Create names array list with
	exchangeList := make([]models.NameEntry, 0)

	// Assign a name to each person
	exchangeList, err = runSecretSanta(persons, exchangeList)

	if err != nil {
		sendError(&err, messagecode.E_SERVER_ERROR, c)
	}

	log.Printf("Secret Santa List Done: %s", exchangeList)

	// Store Secret Santa List to DB
	secretSantaModel := &models.SecretSanta{}
	secretSantaCollection := mongoStore.C(secretSantaModel.Collection())

	secretSantaModel.Uid = models.NewUid()
	secretSantaModel.CreatedAt = time.Now().Unix()
	secretSantaModel.List = exchangeList

	err = secretSantaCollection.Insert(secretSantaModel)
	if err != nil {
		sendError(&err, messagecode.E_SERVER_ERROR, c)
		return
	}

	err = secretSantaCollection.EnsureIndex(secretSantaModel.Index())
	if err != nil {
		sendError(&err, messagecode.E_SERVER_ERROR, c)
		return
	}

	sendResponse(secretSantaModel.Uid, messagecode.S_RESOURCE_OK, c)
	return

}

func runSecretSanta(persons []models.Person, secretSantaList []models.NameEntry) ([]models.NameEntry, error) {
	// For every person, assign them a secret santa

	if len(persons) == 0 {
		return nil, errors.New("I don't any Santa Clauses :(")
	}

	rand.Seed(time.Now().UTC().UnixNano())

	for len(secretSantaList) < len(persons) {
		bag := make([]models.Person, len(persons))
		copy(bag, persons)
		reset := false
		secretSantaList = secretSantaList[:0]
		log.Printf("Selecting Names, names selected (should be 0): %v", len(secretSantaList))
		log.Printf("Bag of names: %v", bag)
		log.Printf("Persons: %v", persons)
		for _, person := range persons {
			log.Printf("%s, it's your turn!", person.Name)
			// Loop until unused name is chosen that is not current person
			var i int

			for i = randInt(0, len(bag)); bag[i].Name == person.Name; i = (i + 1) % len(bag) {
				log.Printf("Chose again %s", person.Name)
				// log.Printf("current index: %v", i)
				// log.Printf("(%v == %v)", bag[i].Name, person.Name)
				// log.Printf("Total Names Chosen: %v", len(secretSantaList))
				// log.Printf("Name Selected: %s", bag[i].Name)
				// log.Printf("Chooser: %s", person.Name)
				// log.Printf("List: %v", secretSantaList)
				// log.Printf("Bag: %v", bag)

				// Edge case that we are on last name and is same person, restart!
				if len(secretSantaList) == len(persons)-1 {
					log.Printf("Last name was equal to selector, restarting Secret santa.")
					reset = true
					break
				}
			}

			if reset {
				break
			}

			// Selected secret santa, update secretSanta list with secret santa
			log.Printf("Name selected for %s", person.Name)
			secretSantaName := models.NameEntry{}
			secretSantaName.Name = bag[i].Name
			secretSantaName.Owner = person
			secretSantaName.AssignedOn = time.Now().Unix()
			secretSantaList = append(secretSantaList, secretSantaName)
			// log.Printf("List: %v", secretSantaList)

			// remove selected name from bag
			bag = append(bag[:i], bag[i+1:]...)
		}
	}
	return secretSantaList, nil
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func spliceToInterface(array interface{}) []interface{} {

	v := reflect.ValueOf(array)
	t := v.Type()

	if t.Kind() != reflect.Slice {
		log.Panicf("`array` should be %s but got %s", reflect.Slice, t.Kind())
	}

	result := make([]interface{}, v.Len(), v.Len())

	for i := 0; i < v.Len(); i++ {
		result[i] = v.Index(i).Interface()
	}

	return result
}
