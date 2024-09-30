package routes

import (
	"slack-chatbot/controller"
	"slack-chatbot/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(middleware.AuthorizeJWT())
	r.POST("/generate", controller.GenerateHandler)
	r.POST("/get-response", controller.GetResponseByPrompt)

}
