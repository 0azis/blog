package http

import (
	"blog/internal/adapter/primary/http/controller"
	"blog/internal/adapter/primary/http/middleware"
	"blog/internal/adapter/secondary/store"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine, store store.Store, savePath string) {
	api := r.Group("/api")
	v1 := api.Group("/v1")
	v1.Static("/uploads", savePath)

	userRoutes(v1, store)
	postRoutes(v1, store)
	relationRoutes(v1, store)
	tagRoutes(v1, store)
	imageRoutes(v1, savePath)
	commentRoutes(v1, store)
}

func userRoutes(r *gin.RouterGroup, store store.Store) {
	controllers := controller.NewUserControllers(store)

	user := r.Group("/users")
	user.GET(":username", middleware.AuthMiddleware, controllers.GetByUsername)
	user.GET("/account", middleware.AuthMiddleware, controllers.Profile)
	user.GET("/search", middleware.AuthMiddleware, controllers.Search)
	user.PATCH("", middleware.AuthMiddleware, controllers.UpdateAccount)

	auth := user.Group("/auth")
	auth.POST("/signin", controllers.SignIn)
	auth.POST("/signup", controllers.SignUp)
	auth.POST("/refresh", middleware.RefreshMiddleware, controllers.RefreshTokens)
	auth.POST("/logout", middleware.AuthMiddleware, controllers.Logout)

	// test
	auth.POST("/test/signin", controllers.SignShitIn)
}

func postRoutes(r *gin.RouterGroup, store store.Store) {
	post := r.Group("/posts", middleware.AuthMiddleware)
	draft := r.Group("/drafts")
	controllers := controller.NewPostControllers(store)

	draft.POST("", controllers.Create)
	draft.POST(":id", controllers.Publish)
	draft.PATCH(":id", controllers.UpdatePost)
	draft.GET("", controllers.GetDrafts)
	draft.GET(":id", controllers.GetDraft)

	post.GET("/author/:id", controllers.GetPostsByUser)
	post.GET(":id", controllers.GetByID)
	post.GET("/author", controllers.MyPosts)
	post.GET("", controllers.GetPosts)
}

func relationRoutes(r *gin.RouterGroup, store store.Store) {
	relation := r.Group("/relation", middleware.AuthMiddleware)
	controllers := controller.NewRelationControllers(store)

	relation.POST("/subscribers/:id", controllers.Subscribe)
	relation.GET("/subscribers/:id", controllers.Subscribers)
	relation.GET("/followers/:id", controllers.Followers)
}

func tagRoutes(r *gin.RouterGroup, store store.Store) {
	tag := r.Group("/tags", middleware.AuthMiddleware)
	controllers := controller.NewTagControllers(store)

	tag.PATCH("", controllers.Create)
	tag.GET("/post/:id", controllers.GetByPostID)
	tag.GET("/top", controllers.GetByPopularity)
}

func imageRoutes(r *gin.RouterGroup, savePath string) {
	image := r.Group("/uploads", middleware.AuthMiddleware)
	controllers := controller.NewImageControllers(savePath)

	image.POST("", controllers.Upload)
}

func commentRoutes(r *gin.RouterGroup, store store.Store) {
	comment := r.Group("/comments", middleware.AuthMiddleware)
	controllers := controller.NewCommentControllesr(store)

	comment.POST("", controllers.NewComment)
	comment.GET("/post/:id", controllers.GetCommentsByPost)
	comment.GET(":id", controllers.GetComment)
}
