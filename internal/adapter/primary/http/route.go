package http

import (
	"blog/internal/adapter/primary/http/controller"
	"blog/internal/adapter/primary/http/middleware"
	"blog/internal/adapter/secondary/store"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine, store store.Store) {
	api := r.Group("/api")
	v1 := api.Group("/v1")

	userRoutes(v1, store)
	postRoutes(v1, store)
	relationRoutes(v1, store)
}

func userRoutes(r *gin.RouterGroup, store store.Store) {
	user := r.Group("/users")
	controllers := controller.NewUserControllers(store)

	user.POST("/signin", controllers.SignIn)
	user.POST("/signup", controllers.SignUp)
	user.GET("/profile", middleware.AuthMiddleware, controllers.Profile)
	user.GET("/search", middleware.AuthMiddleware, controllers.Search)
	user.POST("/refresh", middleware.RefreshMiddleware, controllers.RefreshTokens)
}

func postRoutes(r *gin.RouterGroup, store store.Store) {
	post := r.Group("/posts", middleware.AuthMiddleware)
	controllers := controller.NewPostControllers(store)

	post.POST("/", controllers.Create)
	post.POST("/publish/:id", controllers.Publish)
	post.GET("/", controllers.GetAll)
	post.GET("/:id", controllers.GetOne)
	post.PATCH("/:id", controllers.UpdatePost)
}

func relationRoutes(r *gin.RouterGroup, store store.Store) {
	relation := r.Group("/relation", middleware.AuthMiddleware)
	controllers := controller.NewRelationControllers(store)

	relation.POST("/subscribers/:id", controllers.Subscribe)
	relation.GET("/subscribers", controllers.Subscribers)
	relation.GET("/followers", controllers.Followers)
}
