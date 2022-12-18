package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/Bryan-an/tasker-backend/pkg/common/models"
	"github.com/Bryan-an/tasker-backend/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
)

type ResendCodeInput struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

func (h handler) ResendCode(c *gin.Context) {
	var input ResendCodeInput

	if err := c.ShouldBindJSON(&input); err != nil {
		var ve validator.ValidationErrors

		if errors.As(err, &ve) {
			out := utils.FillErrors(ve)

			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
		} else {
			c.AbortWithError(http.StatusBadRequest, err)
		}

		return
	}

	coll := h.DB.Collection("verifications")
	verificationsFilter := bson.D{{Key: "email", Value: input.Email}}
	result, err := coll.DeleteOne(context.TODO(), verificationsFilter)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if result.DeletedCount == 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": "unable to replace email verification code"})
		return
	}

	u := models.User{
		Name:  input.Name,
		Email: input.Email,
	}

	err = SendVerificationEmail(h, c, u)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "please check your email for email verification code",
	})
}
