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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type resendCodeInput struct {
	Email *string `json:"email" binding:"required,email"`
}

func (h handler) ResendCode(c *gin.Context) {
	var input resendCodeInput

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

	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)
	session, err := h.Client.StartSession()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	defer session.EndSession(context.TODO())

	_, err = session.WithTransaction(
		context.TODO(),
		func(ctx mongo.SessionContext) (interface{}, error) {
			coll := h.DB.Collection("verifications")
			verificationsFilter := bson.D{{Key: "email", Value: input.Email}}
			result, err := coll.DeleteOne(ctx, verificationsFilter)

			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return nil, err
			}

			if result.DeletedCount == 0 {
				c.AbortWithStatusJSON(http.StatusInternalServerError,
					gin.H{"error": "unable to replace email verification code"})

				return nil, errors.New("unable to replace email verification code")
			}

			u := models.User{
				Email: input.Email,
			}

			err = SendVerificationEmail(h, c, u, ctx)

			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return nil, err
			}

			return nil, nil
		}, txnOptions)

	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "please check your email for email verification code",
	})
}
