package config

import (
	"fmt"
	"os"
)

type config struct {
	Server *httpServer
	Db     *database
}

func NewConfig() *config {
	http := initHtppServer()
	db := initDatabase()
	return &config{
		http, db,
	}
}

type httpServer struct {
	host     string
	port     string
	SavePath string
}

func (hs httpServer) BuildSocket() string {
	return hs.host + ":" + hs.port
}

func initHtppServer() *httpServer {
	return &httpServer{
		host:     getEnv("HTTP_HOST", "localhost"),
		port:     getEnv("HTTP_PORT", "8000"),
		SavePath: getEnv("SAVE_IMAGE_PATH", ""),
	}
}

type database struct {
	nameDB   string
	user     string
	password string
	host     string
	port     string
}

func (d database) BuildURI() string {
	uri := fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true", d.user, d.password, d.host, d.port, d.nameDB)
	return uri
}

func initDatabase() *database {
	return &database{
		nameDB:   getEnv("DATABASE_NAME", ""),
		user:     getEnv("DATABASE_USER", ""),
		password: getEnv("DATABASE_PASSWORD", ""),
		host:     getEnv("DATABASE_HOST", "localhot"),
		port:     getEnv("DATABSE_PORT", "3306"),
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
