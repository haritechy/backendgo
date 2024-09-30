package controller

import (
	"log"
	"net/http"
	"os"
	"slack-chatbot/database"
	"slack-chatbot/models"
	"slack-chatbot/requsts"
	"strings"

	"github.com/gin-gonic/gin"
)

func GenerateHandler(c *gin.Context) {
	var prompt models.Prompt
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if err := c.BindJSON(&prompt); err != nil {
		log.Printf("Error binding request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	responseText, err := requsts.GetGeminiResponse(geminiAPIKey, prompt.Prompt)
	if err != nil {
		log.Printf("Error getting response from Gemini: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't process your request"})
		return
	}
	promptResponse := models.PromptResponse{
		Prompt:   strings.TrimSpace(prompt.Prompt),
		Response: strings.TrimSpace(responseText),
	}

	if err := database.DB.Create(&promptResponse).Error; err != nil {
		log.Printf("Error saving to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't save the response"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": responseText})
}
