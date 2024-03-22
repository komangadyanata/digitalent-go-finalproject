package routes

import (
	"mygram/controllers"
	"mygram/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	userRouter := r.Group("/users")
	{
		userRouter.POST("/register", controllers.UserRegister)
		userRouter.POST("/login", controllers.UserLogin)
		userRouter.PUT("/", middlewares.Authentication(), controllers.UpdateUser)
		userRouter.DELETE("/", middlewares.Authentication(), controllers.DeleteUser)
	}
}
