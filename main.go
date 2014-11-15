package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iToto/jollyHelper/common"
	"github.com/iToto/jollyHelper/resources"
	"log"
	"os"
	"runtime"
)

var (
	APP_PORT    = os.Getenv("PORT")
	APP_ENV     = os.Getenv("ENV")
	APP_DB_URL  = os.Getenv("DB_URL")
	APP_DB_NAME = os.Getenv("DB_NAME")
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Printf("APP_ENV - %s = %s", "ENV", APP_ENV)
	log.Printf("APP_PORT  - %s = %s", "PORT", APP_PORT)
	log.Printf("APP_DB_URL  - %s = %s", "DB_URL", APP_DB_URL)
	log.Printf("APP_DB_NAME  - %s = %s", "DB", APP_DB_NAME)

	if APP_ENV == "" || APP_PORT == "" || APP_DB_URL == "" || APP_DB_NAME == "" {
		log.Printf("Missing environment variables: ENV: %s, PORT: %s, DB_URL: %s, DB_NAME  \n", APP_ENV, APP_PORT, APP_DB_URL, APP_DB_NAME)
	}
}

func main() {
	router := gin.Default()

	// Connect to DB
	router.Use(common.MongoDbHandler(APP_DB_URL, APP_DB_NAME))

	personResource := &resources.PersonResource{}
	person := router.Group("/persons/")
	person.POST("/", personResource.Create)
	person.GET("/:id", personResource.Get)
	person.GET("/", personResource.List)
	// person.PUT("/:uid", personResource.Update)
	// person.DELETE("/:uid/:disable", personResource.Disable)

	router.GET("/", func(c *gin.Context) {
		c.String(200, "hello world")
	})
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	router.POST("/submit", func(c *gin.Context) {
		c.String(401, "not authorized")
	})
	router.PUT("/error", func(c *gin.Context) {
		c.String(500, "and error hapenned :(")
	})
	router.Run(":" + APP_PORT)
}
