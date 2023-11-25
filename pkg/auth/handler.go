package auth

import (
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

	routes := r.Group("/api/v1/auth")
	routes.POST("/register", h.Register)
	routes.POST("/login", h.Login)
	routes.GET("/login/facebook", h.InitFacebookLogin)
	routes.POST("/login/facebook/mobile", h.LoginWithFacebookMobile)
	routes.GET("/facebook/callback", h.HandleFacebookLogin)
	routes.GET("/login/google", h.InitGoogleLogin)
	routes.POST("/login/google/mobile", h.LoginWithGoogleMobile)
	routes.GET("/google/callback", h.HandleGoogleLogin)
	routes.POST("/verify/email", h.VerifyEmail)
	routes.POST("/verify/resendCode", h.ResendCode)
}
