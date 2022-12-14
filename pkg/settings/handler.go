package settings

import (
	"github.com/Bryan-an/tasker-backend/pkg/common/middlewares"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type handler struct {
	DB *mongo.Database
}

func RegisterRoutes(r *gin.Engine, db *mongo.Database) {
	h := &handler{
		DB: db,
	}

	routes := r.Group("/api/v1/settings")

	routes.Use(middlewares.JwtAuthMiddleware())
	routes.GET("/", h.GetSettings)
	routes.POST("/", h.AddSettings)
	routes.PATCH("/", h.UpdateSettings)
	routes.PUT("/", h.ReplaceSettings)
}
