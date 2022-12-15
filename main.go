package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Bryan-an/tasker-backend/pkg/auth"
	"github.com/Bryan-an/tasker-backend/pkg/common/db"
	"github.com/Bryan-an/tasker-backend/pkg/common/middlewares"
	"github.com/Bryan-an/tasker-backend/pkg/settings"
	"github.com/Bryan-an/tasker-backend/pkg/tasks"
	"github.com/Bryan-an/tasker-backend/pkg/users"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var database *mongo.Database

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error laading .env file")
		return
	}

	client := db.Connect()
	DbName := os.Getenv("DB_NAME")
	database = client.Database(DbName)

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	router := setupRouter()
	var port string

	if port = os.Getenv("PORT"); port == "" {
		port = ":8080"
	}

	router.Run(port)
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(cors.Default())
	router.Use(middlewares.ErrorHandler())
	router.SetTrustedProxies(nil)

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	auth.RegisterRoutes(router, database)
	settings.RegisterRoutes(router, database)
	tasks.RegisterRoutes(router, database)
	users.RegisterRoutes(router, database)

	return router
}
