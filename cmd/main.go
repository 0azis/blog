package main

import (
	"blog/internal/adapter/primary/http"
	"blog/internal/adapter/secondary/store"
	"blog/internal/config"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	_ "blog/cmd/docs"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		slog.Error("environment not found...")
	}

	config := config.NewConfig()
	httpSocket := config.Server.BuildSocket()
	databaseURI := config.Db.BuildURI()

	store, err := store.NewStore(databaseURI)
	if err != nil {
		slog.Error("error while connecting to database...")
	}

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))	
	http.InitRoutes(r, store)
	err = r.Run(httpSocket)
	if err != nil {
		slog.Error("error while creating a web server...")
	}
}
