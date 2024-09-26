package http

import (
	"blog/internal/adapter/primary/http/controller"
	"blog/internal/adapter/secondary/store"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine, store *store.Store) {
	api := r.Group("/api")
	v1 := api.Group("/v1")

	userRoutes(v1, store)
}

func userRoutes(r *gin.RouterGroup, store *store.Store) {
	controllers := controller.NewUserControllers(store)

	r.POST("/signin", controllers.SignIn)
	r.POST("/signup", controllers.SignUp)
}
