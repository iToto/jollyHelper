package common

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"log"
	//"gopkg.in/mgo.v2/bson"
)

// MongoDBHandler returns a map of the gin context and the MongoDB instance
func MongoDbHandler(url string, database string) gin.HandlerFunc {
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	log.Printf("Successfully connected to Mongo %s:%s", url, database)
	return func(c *gin.Context) {
		c.Set("mongoSession", session)

		s2 := session.Clone()
		c.Set("mongoStore", s2.DB(database))

		defer s2.Close()

		c.Next()
	}
}
