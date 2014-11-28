package resources

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	// "github.com/iToto/jollyHelper/common"
	"github.com/iToto/jollyHelper/common/messagecode"
	"github.com/iToto/jollyHelper/models"
	"log"
	// "github.com/op/go-logging"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// "strconv"
	"time"
)

type PersonResource struct {
}

func (p *PersonResource) Create(c *gin.Context) {
	person := &models.Person{}

	err := binding.JSON.Bind(c.Request, person)

	if err != nil {
		sendError(&err, messagecode.E_INVALID_REQUEST, c)
		return
	}

	if person.Uid == "" {
		person.Uid = models.NewUid()
	}

	mongoStore := c.MustGet("mongoStore").(*mgo.Database)
	personsCollection := mongoStore.C(person.Collection())
	err = personsCollection.EnsureIndex(person.Index())
	if err != nil {
		sendError(&err, messagecode.E_SERVER_ERROR, c)
		return
	}

	query := bson.M{"email": person.Email}
	person.CreatedAt = time.Now().Unix()
	update := bson.M{"$setOnInsert": models.Struct2Map(person)}

	info, err := personsCollection.Upsert(query, update)
	if err != nil {
		log.Print("ERROR: %s", info)
		sendError(&err, messagecode.E_SERVER_ERROR, c)
		return
	}

	sendResponse(&person, messagecode.S_RESOURCE_CREATED, c)
	return
}

func (p *PersonResource) Get(c *gin.Context) {

	id := c.Params.ByName("id")

	if id == "" {
		err := errors.New("invalid id")
		sendError(&err, messagecode.E_INVALID_REQUEST, c)
		return
	}

	log.Printf("Reading Person with UID: %s", id)

	var err error
	mongoStore := c.MustGet("mongoStore").(*mgo.Database)

	person := &models.Person{}
	personsCollection := mongoStore.C(person.Collection())

	err = personsCollection.EnsureIndex(person.Index())
	if err != nil {
		sendError(&err, messagecode.E_SERVER_ERROR, c)
		return
	}

	err = personsCollection.Find(bson.M{"uid": id}).One(person)

	log.Printf("Person %s", person)

	if err != nil {
		if err == mgo.ErrNotFound {
			sendError(&err, messagecode.S_RESOURCE_NOTFOUND, c)
		} else {
			sendError(&err, messagecode.E_SERVER_ERROR, c)
		}
		return
	}
	sendResponse(&person, messagecode.S_RESOURCE_OK, c)
	return
}

func (p *PersonResource) List(c *gin.Context) {
	var err error
	mongoStore := c.MustGet("mongoStore").(*mgo.Database)

	person := &models.Person{}
	persons := []models.Person{}
	personCollection := mongoStore.C(person.Collection())

	err = personCollection.EnsureIndex(person.Index())

	if err != nil {
		sendError(&err, messagecode.E_SERVER_ERROR, c)
		return
	}

	err = personCollection.Find(bson.M{}).Limit(person.Limit()).All(&persons)

	log.Printf("persons: %q", person)

	if err != nil {
		if err == mgo.ErrNotFound {
			sendError(&err, messagecode.S_RESOURCE_NOTFOUND, c)
		} else {
			sendError(&err, messagecode.E_SERVER_ERROR, c)
		}
		return
	}
	sendResponse(&persons, messagecode.S_RESOURCE_OK, c)
	return
}

func (p *PersonResource) AddListItem(c *gin.Context) {
	id := c.Params.ByName("id")

	if id == "" {
		err := errors.New("invalid id")
		sendError(&err, messagecode.E_INVALID_REQUEST, c)
		return
	}

	var err error
	mongoStore := c.MustGet("mongoStore").(*mgo.Database)

	person := &models.Person{}
	personsCollection := mongoStore.C(person.Collection())

	err = personsCollection.EnsureIndex(person.Index())
	if err != nil {
		sendError(&err, messagecode.E_SERVER_ERROR, c)
		return
	}

	err = personsCollection.Find(bson.M{"uid": id}).One(person)

	log.Printf("Person %s", person)

	if err != nil {
		if err == mgo.ErrNotFound {
			sendError(&err, messagecode.S_RESOURCE_NOTFOUND, c)
		} else {
			sendError(&err, messagecode.E_SERVER_ERROR, c)
		}
		return
	}

	// Add new list item
	listItem := &models.ListItem{}

	err = binding.JSON.Bind(c.Request, listItem)

	if err != nil {
		sendError(&err, messagecode.E_INVALID_REQUEST, c)
		return
	}

	if listItem.Uid == "" {
		listItem.Uid = models.NewUid()
	}

	log.Printf("List received: %v", listItem)

	// person.List = append(person.List, listItem)
	person.List = append(person.List, *listItem)

	log.Printf("Person to be updated: %v", person)

	query := bson.M{"email": person.Email}
	person.UpdatedAt = time.Now().Unix()
	update := bson.M{"$set": models.Struct2Map(person)}

	err = personsCollection.Update(query, update)
	if err != nil {
		log.Print("ERROR: %s", err)
		sendError(&err, messagecode.E_SERVER_ERROR, c)
		return
	}

	sendResponse(&person, messagecode.S_RESOURCE_CREATED, c)
	return

}
