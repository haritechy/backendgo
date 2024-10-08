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

	// Retrieve the user ID from the context
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint, ok := userID.(int) // Ensure userID is of the correct type (int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID type"})
		return
	}

	// Bind the prompt from the request body
	if err := c.BindJSON(&prompt); err != nil {
		log.Printf("Error binding request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Get the response from Gemini API
	responseText, err := requsts.GetGeminiResponse(geminiAPIKey, prompt.Prompt)
	if err != nil {
		log.Printf("Error getting response from Gemini: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't process your request"})
		return
	}

	// Create the PromptResponse object
	promptResponse := models.PromptResponse{
		Prompt:   strings.TrimSpace(prompt.Prompt),
		Response: strings.TrimSpace(responseText),
		UserID:   userIDUint, // Store the user ID here
	}

	// Save the prompt and response in the database
	if err := database.DB.Create(&promptResponse).Error; err != nil {
		log.Printf("Error saving to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't save the response"})
		return
	}

	// Return the response back to the client
	c.JSON(http.StatusOK, gin.H{"response": responseText})
}

func GetHistoryHandler(c *gin.Context) {
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var promptResponses []models.PromptResponse

	// Fetching prompt responses for the user
	if err := database.DB.Where("user_id = ?", userID).Find(&promptResponses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't fetch history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"history": promptResponses})
}
