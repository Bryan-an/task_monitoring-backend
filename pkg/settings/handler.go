package settings

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

	routes := r.Group("/api/v1/settings")

	routes.Use(middlewares.JwtAuthMiddleware())
	routes.GET("/", h.GetSettings)
	routes.POST("/", h.AddSettings)
	routes.PUT("/", h.ReplaceSettings)
	routes.PATCH("/", h.UpdateSettings)
}
