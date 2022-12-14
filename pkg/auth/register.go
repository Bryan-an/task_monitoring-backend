package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Bryan-an/tasker-backend/pkg/common/models"
	"github.com/Bryan-an/tasker-backend/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h handler) Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		var ve validator.ValidationErrors

		if errors.As(err, &ve) {
			out := make([]utils.ErrorMsg, len(ve))

			for i, fe := range ve {
				out[i] = utils.ErrorMsg{
					Field:   fe.Field(),
					Message: utils.GetErrorMsg(fe),
				}
			}

			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		return
	}

	usersCollection := h.DB.Collection("users")
	filter := bson.D{{Key: "email", Value: input.Email}}
	var user models.User
	err := usersCollection.FindOne(context.TODO(), filter).Decode(&user)

	if err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "this email address is already in use",
		})
		return
	}

	if input.Name == "" {
		input.Name = input.Email
	}

	hash, err := utils.HashPassword(input.Password)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.Password = hash

	u := models.User{
		Name:      input.Name,
		Email:     input.Email,
		Password:  input.Password,
		Role:      "user",
		Status:    "created",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req, err := usersCollection.InsertOne(context.TODO(), u)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error registering user",
		})
		return
	}

	s := models.Settings{
		UserId: fmt.Sprint(req.InsertedID.(primitive.ObjectID).Hex()),
		Notifications: models.Notification{
			Email:  false,
			Mobile: true,
		},
		Security:  "something",
		Theme:     "light",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	settingsCollection := h.DB.Collection("settings")
	_, err = settingsCollection.InsertOne(context.TODO(), s)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error registering settings",
		})
		log.Fatal(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user registrated successfully",
		"id":      req.InsertedID,
	})
}
