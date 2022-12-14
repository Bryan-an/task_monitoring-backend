package settings

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Bryan-an/tasker-backend/pkg/common/models"
	"github.com/Bryan-an/tasker-backend/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Notification struct {
	Email  bool `json:"email"`
	Mobile bool `json:"mobile"`
}

type SettingsInput struct {
	Notifications Notification `json:"notifications" binding:"required"`
	Security      string       `json:"security" binding:"required"`
	Theme         string       `json:"theme" binding:"required,oneof=dark light"`
}

func (h handler) AddSettings(c *gin.Context) {
	uid, err := utils.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Fatal(err)
		return
	}

	var input SettingsInput

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

	s := models.Settings{
		UserId: uid,
		Notifications: models.Notification{
			Email:  input.Notifications.Email,
			Mobile: input.Notifications.Mobile,
		},
		Security:  input.Security,
		Theme:     input.Theme,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	settingsCollection := h.DB.Collection("settings")
	req, err := settingsCollection.InsertOne(context.TODO(), s)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error registering settings",
		})
		log.Fatal(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "settings added",
		"id":      req.InsertedID,
	})
}
