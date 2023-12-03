package tasks

import (
	"github.com/Bryan-an/tasker-backend/pkg/common/middlewares"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type handler struct {
	DB     *mongo.Database
	Client *mongo.Client
}

func RegisterRoutes(r *gin.Engine, db *mongo.Database, client *mongo.Client) {
	h := &handler{
		DB:     db,
		Client: client,
	}

	routes := r.Group("/api/v1/tasks")

	routes.Use(middlewares.JwtAuthMiddleware())
	routes.GET("/", h.GetTasks)
	routes.GET("/today", h.GetTasksForToday)
	routes.POST("/", h.AddTask)
	routes.GET("/:id", h.GetTask)
	routes.PUT("/:id", h.ReplaceTask)
	routes.PATCH("/:id", h.UpdateTask)
	routes.DELETE("/:id", h.DeleteTask)
}
