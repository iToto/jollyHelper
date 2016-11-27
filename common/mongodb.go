package common

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"log"
	"time"
	//"gopkg.in/mgo.v2/bson"
)

// MongoDBHandler returns a map of the gin context and the MongoDB instance
func MongoDbHandler(hosts, database, username, password string) gin.HandlerFunc {
	connInfo := mgo.DialInfo{
		Addrs:    []string{hosts},
		Timeout:  60 * time.Second,
		Database: database,
		Username: username,
		Password: password,
	}

	log.Printf("Host:%s, DB:%s, Username:%s, PW:%s", hosts, database, username, password)
	log.Printf("Connecting with info: %v", connInfo)
	session, err := mgo.DialWithInfo(&connInfo)
	if err != nil {
		log.Printf("Could not connect to DB: %v", err)
		panic(err)
	}
	log.Print("Successfully connected to Mongo")
	return func(c *gin.Context) {
		c.Set("mongoSession", session)

		s2 := session.Clone()
		c.Set("mongoStore", s2.DB(connInfo.Database))

		defer s2.Close()

		c.Next()
	}
}
