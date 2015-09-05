package resources

import (
	"errors"
	// "github.com/freehaha/token-auth"
	// "github.com/freehaha/token-auth/memory"
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
	"code.google.com/p/go-uuid/uuid"
	"strings"
	"time"
)

type PersonResource struct {
}

// Auth Functions

// Get user by email address and compare password.
// If correct, create a session object and return the token
func (p *PersonResource) Login(c *gin.Context) {
	credentials := &models.Login{}

	c.Bind(&credentials)

	mongoStore := c.MustGet("mongoStore").(*mgo.Database)

	person := &models.Person{}
	personsCollection := mongoStore.C(person.Collection())

	err := personsCollection.EnsureIndex(person.Index())
	if err != nil {
		sendError(&err, messagecode.E_SERVER_ERROR, c)
		return
	}

	err = personsCollection.Find(bson.M{"email": credentials.Username}).One(person)

	bytePassword := []byte(credentials.Password)
	salt := person.PasswordSalt()
	hashedPassword := person.HashPassword(bytePassword, salt)

	if strings.EqualFold(hashedPassword, person.Password) {
		log.Printf("Authentication Successful")

		uuid := uuid.New()
		log.Printf("token: %s", uuid)
		person.Token = uuid

		query := bson.M{
			"email": person.Email,
		}

		updatedPerson := bson.M{"$set": models.Struct2Map(person)}

		err = personsCollection.Update(query, updatedPerson)

		if err != nil {
			sendError(&err, messagecode.E_SERVER_ERROR, c)
			return
		}
		token := &models.Token{}
		token.Token = uuid
		token.TTL = 999
		token.Owner = person.Email

		sendResponse(&token, messagecode.S_RESOURCE_OK, c)

	} else {
		log.Printf("Authentication Failure")
	}

}

func (p *PersonResource) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestToken := c.Request.URL.Query().Get("token")

		log.Printf("Validating Request Token: %s", requestToken)

		mongoStore := c.MustGet("mongoStore").(*mgo.Database)

		person := &models.Person{}
		personsCollection := mongoStore.C(person.Collection())

		err := personsCollection.EnsureIndex(person.Index())
		if err != nil {
			sendError(&err, messagecode.E_SERVER_ERROR, c)
			return
		}

		err = personsCollection.Find(bson.M{"token": requestToken}).One(person)

		if err != nil {
			c.Fail(401, err)
		}
		c.Set("token", requestToken)

	}
}

func getByEmail(c *gin.Context, email string) (found bool, person *models.Person) {
	var err error
	foundPerson := &models.Person{}
	mongoStore := c.MustGet("mongoStore").(*mgo.Database)
	personsCollection := mongoStore.C(foundPerson.Collection())

	err = personsCollection.EnsureIndex(foundPerson.Index())
	if err != nil {
		sendError(&err, messagecode.E_SERVER_ERROR, c)
		return
	}

	log.Printf("Verifying if email exists: %s", email)
	err = personsCollection.Find(bson.M{"email": email}).One(foundPerson)
	log.Printf("Person found %s", foundPerson)
	// Person does not exist
	if foundPerson == nil {
		log.Printf("Email not in use: %s", email)
		return false, nil
	}

	if err != nil {
		if err == mgo.ErrNotFound {
			return false, nil
		} else {
			sendError(&err, messagecode.E_SERVER_ERROR, c)
		}
		return false, nil
	}

	log.Printf("Person found %s", foundPerson)
	return true, foundPerson
}

func createPerson(c *gin.Context, person *models.Person) error {
	if person.Uid == "" {
		person.Uid = models.NewUid()
	}

	bytePassword := []byte(person.Password)
	salt := person.PasswordSalt()
	person.Password = person.HashPassword(bytePassword, salt)

	mongoStore := c.MustGet("mongoStore").(*mgo.Database)
	personsCollection := mongoStore.C(person.Collection())
	err := personsCollection.EnsureIndex(person.Index())

	if err != nil {
		return err
	}

	query := bson.M{"email": person.Email}
	person.CreatedAt = time.Now().Unix()
	update := bson.M{"$setOnInsert": models.Struct2Map(person)}

	info, err := personsCollection.Upsert(query, update)

	if err != nil {
		log.Print("ERROR: %s", info)
		return err
	}

	return nil
}

func (p *PersonResource) Register(c *gin.Context) {
	person := &models.Person{}
	err := binding.JSON.Bind(c.Request, person)

	if err != nil {
		sendError(&err, messagecode.E_INVALID_REQUEST, c)
		return
	}

	// Ensure user does not already exist
	emailExists, _ := getByEmail(c, person.Email)
	if emailExists {
		// Cannot register, email already taken!
		err := errors.New("Email invalid, please try again")
		sendError(&err, messagecode.E_INVALID_REQUEST, c)
		return
	}

	log.Printf("Proceeding to create user: %s", person)

	// Create the person!
	err = createPerson(c, person)

	if err != nil {
		err := errors.New("Could not create person")
		sendError(&err, messagecode.E_SERVER_ERROR, c)
		return
	}

	sendResponse(&person, messagecode.S_RESOURCE_CREATED, c)
	return
}

func (p *PersonResource) Create(c *gin.Context) {
	log.Printf("Call to create a user")
	person := &models.Person{}

	err := binding.JSON.Bind(c.Request, person)

	if err != nil {
		log.Printf("Could not bind input")
		sendError(&err, messagecode.E_INVALID_REQUEST, c)
		return
	}
	err = createPerson(c, person)

	if err != nil {
		err := errors.New("Could not create person")
		sendError(&err, messagecode.E_SERVER_ERROR, c)
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

func (p *PersonResource) GetUserList(c *gin.Context) {
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

	sendResponse(&person.List, messagecode.S_RESOURCE_OK, c)
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
