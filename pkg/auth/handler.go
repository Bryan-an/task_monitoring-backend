package auth

import (
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

	routes := r.Group("/api/v1/auth")
	routes.POST("/register", h.Register)
	routes.POST("/login", h.Login)
	routes.GET("/login/facebook", h.InitFacebookLogin)
	routes.GET("/facebook/callback", h.HandleFacebookLogin)
}
