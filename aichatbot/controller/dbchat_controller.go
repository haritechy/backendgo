package controller

import (
	"log"
	"net/http"

	"slack-chatbot/database"
	"slack-chatbot/models"

	"github.com/gin-gonic/gin"
)

func GetResponseByPrompt(c *gin.Context) {
	var requestPrompt models.Prompt
	var storedResponse models.PromptResponse
	if err := c.BindJSON(&requestPrompt); err != nil {
		log.Printf("Error binding request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := database.DB.Where("prompt = ?", requestPrompt.Prompt).First(&storedResponse).Error; err != nil {
		log.Printf("Prompt not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "prompt not found in database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": storedResponse.Response})
}
