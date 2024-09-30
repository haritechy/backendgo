package routes

import (
	"slack-chatbot/controller"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	userRoutes := r.Group("/user")
	{
		userRoutes.POST("/create", controller.UserRegister)
		userRoutes.GET("/", controller.UserGet)
		userRoutes.DELETE("/del/:id", controller.UseDelete)
		userRoutes.PUT("/eupdate/:id", controller.UserUpdate)
		userRoutes.POST("/login", controller.UserLogin)
		userRoutes.POST("/chanege", controller.ChnagePassword)
		userRoutes.POST("/for", controller.UserForgot)
		userRoutes.GET("/:id", controller.UserGetbyEmail)
	}
}
