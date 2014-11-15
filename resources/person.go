package resources

import (
	// "errors"
	// "fmt"
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
