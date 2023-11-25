package settings

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Bryan-an/tasker-backend/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
)

type JSONNotifications struct {
	Value notifications
	Valid bool
	Set   bool
}

type notifications struct {
	Email  utils.JSONBool `json:"email"`
	Mobile utils.JSONBool `json:"mobile"`
}

type UpdateInput struct {
	Notifications JSONNotifications `json:"notifications"`
	Theme         utils.JSONString  `json:"theme"`
}

func (n *JSONNotifications) UnmarshalJSON(data []byte) error {
	n.Set = true

	if string(data) == "null" {
		n.Valid = false
		return nil
	}

	var temp notifications

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	n.Value = temp
	n.Valid = true
	return nil
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

	if input.Notifications.Set {
		if input.Notifications.Valid {
			if input.Notifications.Value.Email.Set {
				if input.Notifications.Value.Email.Valid {
					data["notifications.email"] = input.Notifications.Value.Email.Value
				} else {
					data["notifications.email"] = nil
				}
			}

			if input.Notifications.Value.Mobile.Set {
				if input.Notifications.Value.Mobile.Valid {
					data["notifications.mobile"] = input.Notifications.Value.Mobile.Value
				} else {
					data["notifications.mobile"] = nil
				}
			}
		} else {
			data["notifications"] = nil
		}
	}

	if input.Theme.Set {
		if input.Theme.Valid {
			data["theme"] = input.Theme.Value
		} else {
			data["theme"] = nil
		}
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
