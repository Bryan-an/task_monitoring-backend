package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Bryan-an/tasker-backend/pkg/common/models"
	"github.com/Bryan-an/tasker-backend/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func (h handler) VerifyEmail(c *gin.Context) {
	var data models.VerificationData

	if err := c.ShouldBindJSON(&data); err != nil {
		var ve validator.ValidationErrors

		if errors.As(err, &ve) {
			out := utils.FillErrors(ve)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
		} else {
			c.AbortWithError(http.StatusBadRequest, err)
		}

		return
	}

	verificationsColl := h.DB.Collection("verifications")
	verificationsFilter := bson.D{{Key: "email", Value: data.Email}}
	var actualData models.VerificationData
	const verificationNotFoundMessage = "verification code not found for user with email '%s'"

	if err := verificationsColl.FindOne(context.TODO(), verificationsFilter).Decode(&actualData); err != nil {
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf(verificationNotFoundMessage, *data.Email),
			})

			return
		}

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	valid, err := verifyData(actualData, data)

	if !valid {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
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
			usersColl := h.DB.Collection("users")

			usersFilter := bson.D{
				{Key: "email", Value: data.Email},
				{Key: "status", Value: "created"},
			}

			update := bson.D{
				{
					Key: "$set",
					Value: bson.D{
						{Key: "status", Value: "active"},
					},
				},
			}

			result, err := usersColl.UpdateOne(ctx, usersFilter, update)

			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return nil, err
			}

			if result.MatchedCount == 0 {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"error": fmt.Sprintf("user not found with email '%s'", *data.Email),
				})

				return nil, fmt.Errorf("user not found with email '%s'", *data.Email)
			}

			verificationsResult, err := verificationsColl.DeleteOne(ctx, verificationsFilter)

			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return nil, err
			}

			if verificationsResult.DeletedCount == 0 {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"error": fmt.Sprintf(verificationNotFoundMessage, *data.Email),
				})

				return nil, fmt.Errorf(verificationNotFoundMessage, *data.Email)
			}

			return nil, nil
		},
		txnOptions)

	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messasge": "email verified successfully",
	})
}

func verifyData(actualData models.VerificationData, data models.VerificationData) (bool, error) {
	if *actualData.Code != *data.Code {
		return false, errors.New("verification code provided is invalid, please look in your email for the code")
	}

	if actualData.ExpiresAt.Before(time.Now()) {
		return false, errors.New("verification code has expired, please try generating a new code")
	}

	return true, nil
}
