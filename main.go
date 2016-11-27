package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iToto/jollyHelper/common"
	"github.com/iToto/jollyHelper/resources"
	"github.com/mostafah/mandrill"
	"log"
	"os"
	"runtime"
)

var (
	APP_PORT           = os.Getenv("PORT")
	APP_ENV            = os.Getenv("ENV")
	APP_MANDRILL_KEY   = os.Getenv("MANDRILL_KEY")
	APP_MONGO_HOSTS    = os.Getenv("MONGO_HOSTS")
	APP_MONGO_DATABASE = os.Getenv("MONGO_DATABASE")
	APP_MONGO_USERNAME = os.Getenv("MONGO_USERNAME")
	APP_MONGO_PASSWORD = os.Getenv("MONGO_PASSWORD")
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
	log.Printf("MONGO_HOSTS  - %s", APP_MONGO_HOSTS)
	log.Printf("MONGO_DATABASE  - %s", APP_MONGO_DATABASE)
	log.Printf("MONGO_USERNAME  - %s", APP_MONGO_USERNAME)
	log.Printf("MONGO_PASSWORD  - %s", APP_MONGO_PASSWORD)

	if APP_ENV == "" || APP_PORT == "" || APP_MONGO_DATABASE == "" || APP_MONGO_HOSTS == "" || APP_MONGO_USERNAME == "" || APP_MONGO_PASSWORD == "" {
		log.Printf("Missing environment variables: ENV: %s, PORT: %s, MONGOLAB_URI: %s, MONGOLAB_NAME  \n", APP_ENV, APP_PORT, APP_MONGO_HOSTS, APP_MONGO_DATABASE)
		panic("Killing app due to missing environment variables")
	}
}

func main() {
	router := gin.Default()

	// Connect to DB
	router.Use(common.MongoDbHandler(APP_MONGO_HOSTS, APP_MONGO_DATABASE, APP_MONGO_USERNAME, APP_MONGO_PASSWORD))

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

	// Auth
	auth := router.Group("/auth")
	auth.POST("login", personResource.Login)
	auth.POST("register", personResource.Register)

	person := router.Group("/persons")
	person.Use(personResource.Authenticate())
	person.POST("", personResource.Create)
	person.GET("/:id", personResource.Get)
	person.GET("", personResource.List)
	person.POST("/:id/list", personResource.AddListItem)
	person.GET("/:id/list", personResource.GetUserList)
	// person.PUT("/:uid", personResource.Update)
	// person.DELETE("/:uid/:disable", personResource.Disable)

	secretSantaResource := resources.SecretSantaResource{}
	secretSanta := router.Group("/secretsanta")
	secretSanta.Use(personResource.Authenticate())
	secretSanta.POST("", secretSantaResource.AssignNames)
	secretSanta.GET("", secretSantaResource.List)

	notificationResource := resources.NotificationResource{}
	notification := router.Group("/notification")
	notification.Use(personResource.Authenticate())
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
