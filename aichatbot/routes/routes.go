package routes

import (
	"slack-chatbot/controller"
	"slack-chatbot/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/generate", controller.GenerateHandler)
	r.POST("/get-response", controller.GetResponseByPrompt)

}
