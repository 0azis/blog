//	@title		Blog API
//	@version	1.0

//	@host		localhost:8000
//	@BasePath	/api/v1

package main

import (
	"blog/internal/adapter/primary/http"
	"blog/internal/adapter/secondary/store"
	"blog/internal/config"
	"flag"
	"log/slog"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "blog/cmd/docs"
)

func main() {
	devArg := flag.Bool("dev", false, "")
	flag.Parse()

	if *devArg {
		if err := godotenv.Load("../dev/.env"); err != nil {
			slog.Error(err.Error())
		}
	} else {
		if err := godotenv.Load("../.env"); err != nil {
			slog.Error(err.Error())
		}
	}

	config := config.NewConfig()
	httpSocket := config.Server.BuildSocket()
	databaseURI := config.Db.BuildURI()

	store, err := store.NewStore(databaseURI)
	if err != nil {
		slog.Error("error while connecting to database...")
	}

	r := gin.Default()
	cfg := cors.Config{
		AllowOrigins:     []string{"localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(cfg))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	http.InitRoutes(r, store)

	err = r.Run(httpSocket)
	if err != nil {
		slog.Error("error while creating a web server...")
	}
}
