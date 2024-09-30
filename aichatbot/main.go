package main

import (
	"log"
	"slack-chatbot/database"
	"slack-chatbot/middleware"
	"slack-chatbot/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := gin.Default()

	r.Use(middleware.CorsMiddleware())
	database.InitDB()
	routes.SetupRoutes(r)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
