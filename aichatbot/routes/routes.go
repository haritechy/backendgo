package routes

import (
	"slack-chatbot/controller"
	"slack-chatbot/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// Define routes
	protected := r.Group("/chat")
	protected.Use(middleware.AuthorizeJWT())
	{
		protected.POST("/generate", controller.GenerateHandler)
		protected.POST("/get-response", controller.GetResponseByPrompt)
		protected.GET("/history", controller.GetHistoryHandler)

	}

}
