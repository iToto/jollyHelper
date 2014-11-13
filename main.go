package main

import (
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"log"
	"os"
	"runtime"
)

var (
	APP_PORT   = os.Getenv("PORT")
	APP_ENV    = os.Getenv("ENV")
	LOG_FORMAT = os.Getenv("LOG_FORMAT")
	LOGGER     *logging.Logger
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Printf("APP_ENV - %s = %s", "ENV", APP_ENV)
	log.Printf("APP_PORT  - %s = %s", "PORT", APP_PORT)

	if APP_ENV == "" || APP_PORT == "" {
		log.Printf("Please define application environment and port with global variables ENV: %s (local,testing,staging,production) and PORT: %s respectively \n", APP_ENV, APP_PORT)
		APP_ENV = "production" // Default to production if not set...
	}
}

func main() {

	LOGGER = logging.MustGetLogger("jollyHelper")
	logBackend := logging.NewLogBackend(os.Stdout, "", 0)
	logging.SetFormatter(logging.MustStringFormatter(LOG_FORMAT))
	logging.SetBackend(logBackend)

	router := gin.Default()
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
