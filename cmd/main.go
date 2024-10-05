//	@title		Blog API
//	@version	1.0

//	@host		localhost:8000
//	@BasePath	/api/v1

package main

import (
	"blog/internal/adapter/primary/http"
	"blog/internal/adapter/secondary/store"
	"blog/internal/config"
	"fmt"
	"log/slog"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "blog/cmd/docs"
)

func main() {
	if err := godotenv.Load("../.envprod"); err != nil {
		fmt.Println(err)
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

	// setup cors settings
	cfg := cors.DefaultConfig()
	cfg.AllowAllOrigins = true
	cfg.AllowCredentials = true
	r.Use(cors.New(cfg))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	http.InitRoutes(r, store)

	err = r.Run(httpSocket)
	if err != nil {
		slog.Error("error while creating a web server...")
	}
}
