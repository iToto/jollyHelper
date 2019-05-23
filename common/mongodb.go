package common

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"log"
	"time"
	//"gopkg.in/mgo.v2/bson"
)

// MongoDBHandler returns a map of the gin context and the MongoDB instance
func MongoDbHandler(mongoUrl string) gin.HandlerFunc {
	dialInfo, err := mgo.ParseURL(mongoUrl)
	if err != nil {
		log.Fatalf("Could not parse conn url: %s\n", mongoUrl)
	}

	log.Printf("Creating Session: %v", connInfo)
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}

	// log.Printf("Authenticating Session: %v", dbCreds)
	// err = session.Login(dbCreds)
	// if err != nil {
	// 	log.Fatalf("Could not authenticate: %s\n", err)
	// }

	log.Print("Successfully connected to Mongo")
	return func(c *gin.Context) {
		c.Set("mongoSession", session)

		s2 := session.Clone()
		c.Set("mongoStore", s2.DB(connInfo.Database))

		defer s2.Close()

		c.Next()
	}
}
