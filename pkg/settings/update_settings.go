package settings

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Bryan-an/tasker-backend/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
)

type notifications struct {
	Email  *bool `json:"email"`
	Mobile *bool `json:"mobile"`
}

type UpdateInput struct {
	Notifications *notifications `json:"notifications"`
	Theme         *string        `json:"theme"`
}

func (h handler) UpdateSettings(c *gin.Context) {
	uid, err := utils.ExtractTokenID(c)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var input UpdateInput

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

	settingsCollection := h.DB.Collection("settings")
	filter := bson.D{{Key: "user_id", Value: uid}}

	data := bson.M{
		"updated_at": time.Now(),
	}

	if input.Notifications != nil {
		if input.Notifications.Email != nil {
			data["notifications.email"] = input.Notifications.Email
		}

		if input.Notifications.Mobile != nil {
			data["notifications.mobile"] = input.Notifications.Mobile
		}
	}

	if input.Theme != nil {
		data["theme"] = input.Theme
	}

	update := bson.D{
		{
			Key:   "$set",
			Value: data,
		},
	}

	result, err := settingsCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if result.MatchedCount == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf(
				"settings not found for user with id '%s'",
				uid.Hex(),
			),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "settings updated successfully",
	})
}
