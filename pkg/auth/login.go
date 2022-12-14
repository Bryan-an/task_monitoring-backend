package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Bryan-an/tasker-backend/pkg/common/models"
	"github.com/Bryan-an/tasker-backend/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h handler) Login(c *gin.Context) {
	var input LoginInput

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
	var u models.User
	err := usersCollection.FindOne(context.TODO(), filter).Decode(&u)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("user not found with email '%s'", input.Email),
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Fatal(err)
		return
	}

	correct := verifyPassword(input.Password, u.Password)

	if !correct {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "user or password incorrect",
		})
		return
	}

	token, err := utils.GenerateToken(u.Id.Hex())

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Fatal(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
