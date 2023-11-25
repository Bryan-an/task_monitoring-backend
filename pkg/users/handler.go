package users

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

	routes := r.Group("/api/v1/users")

	routes.Use(middlewares.JwtAuthMiddleware())
	routes.GET("/", h.GetUser)
	routes.PUT("/", h.ReplaceUser)
	routes.PATCH("/", h.UpdateUser)
	routes.DELETE("/", h.DeleteUser)
}
