package main

import (
	"github.com/freehaha/token-auth"
	"github.com/freehaha/token-auth/memory"
	"github.com/gin-gonic/gin"
	"github.com/iToto/jollyHelper/common"
	"github.com/iToto/jollyHelper/resources"
	"github.com/mostafah/mandrill"
	"log"
	"os"
	"runtime"
)

var (
	APP_PORT         = os.Getenv("PORT")
	APP_ENV          = os.Getenv("ENV")
	APP_DB_URL       = os.Getenv("MONGOLAB_URI")
	APP_DB_NAME      = os.Getenv("MONGOLAB_NAME")
	APP_MANDRILL_KEY = os.Getenv("MANDRILL_KEY")
)

func init() {

	log.Printf("  ,--.       ,--.,--.         ,--.  ,--.       ,--.                      ")
	log.Printf("  `--' ,---. |  ||  |,--. ,--.|  '--'  | ,---. |  | ,---.  ,---. ,--.--. ")
	log.Printf("  ,--.| .-. ||  ||  | \\  '  / |  .--.  || .-. :|  || .-. || .-. :|  .--' ")
	log.Printf("  |  |' '-' '|  ||  |  \\   '  |  |  |  |\\   --.|  || '-' '\\   --.|  |    ")
	log.Printf(".-'  / `---' `--'`--'.-'  /   `--'  `--' `----'`--'|  |-'  `----'`--'    ")
	log.Printf("'---'                `---'                         `--'                  ")

	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Printf("APP_ENV - %s = %s", "ENV", APP_ENV)
	log.Printf("APP_PORT  - %s = %s", "PORT", APP_PORT)
	log.Printf("APP_DB_URL  - %s = %s", "MONGOLAB_URI", APP_DB_URL)
	log.Printf("APP_DB_NAME  - %s = %s", "MONGOLAB_NAME", APP_DB_NAME)

	if APP_ENV == "" || APP_PORT == "" || APP_DB_URL == "" || APP_DB_NAME == "" || APP_DB_URL == "" {
		log.Printf("Missing environment variables: ENV: %s, PORT: %s, MONGOLAB_URI: %s, MONGOLAB_NAME  \n", APP_ENV, APP_PORT, APP_DB_URL, APP_DB_NAME)
		panic("Killing app due to missing environment variables")
	}
}

func tokenMemoryStore(salt string) gin.HandlerFunc {
	return func(c *gin.Context) {
		memStore := memstore.New(salt)
		tokenAuth := tauth.NewTokenAuth(nil, nil, memStore, nil)

		c.Set("tokenStore", memStore)
		c.Set("tokenAuth", tokenAuth)
	}
}

func main() {
	router := gin.Default()

	// Connect to DB
	router.Use(common.MongoDbHandler(APP_DB_URL, APP_DB_NAME))

	// Setup Token Storage
	router.Use(tokenMemoryStore("jollyHelper"))

	// Test Mandrill
	mandrill.Key = APP_MANDRILL_KEY
	// you can test your API key with Ping
	err := mandrill.Ping()
	// everything is OK if err is nil

	if err != nil {
		log.Printf("Failed to ping Mandrill: %s", err)
		panic("Unable to ping Mandrill")
	}

	personResource := &resources.PersonResource{}
	person := router.Group("/persons")
	person.POST("", personResource.Create)
	person.GET("/:id", personResource.Get)
	person.GET("", personResource.List)
	person.POST("/:id/list", personResource.AddListItem)
	person.GET("/:id/list", personResource.GetUserList)
	// person.PUT("/:uid", personResource.Update)
	// person.DELETE("/:uid/:disable", personResource.Disable)

	// Auth
	auth := router.Group("/auth")
	auth.POST("login", personResource.Login)

	secretSantaResource := resources.SecretSantaResource{}
	secretSanta := router.Group("/secretsanta")
	secretSanta.POST("", secretSantaResource.AssignNames)
	secretSanta.GET("", secretSantaResource.List)

	notificationResource := resources.NotificationResource{}
	notification := router.Group("/notification")
	notification.GET("/:id", notificationResource.Send)

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
