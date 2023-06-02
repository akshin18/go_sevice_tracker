package main

import (
	handler "go_server/handlers"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	var err = godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	router := gin.Default()
	// same as
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.GET("/get_events", handler.GetEvent)
	router.GET("/get_unic_events", handler.GetUnicEventType)
	router.GET("/get_unic_events_name", handler.GetUnicEventName)
	router.GET("/get_unic_browsers", handler.GetUnicBrowserType)
	router.POST("/add_event", handler.PostEvent)

	router.Run("0.0.0.0:8080")
}
