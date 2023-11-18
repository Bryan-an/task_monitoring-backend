package settings

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Bryan-an/tasker-backend/pkg/common/models"
	"github.com/Bryan-an/tasker-backend/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type notification struct {
	Email  *bool `json:"email" binding:"required"`
	Mobile *bool `json:"mobile" binding:"required"`
}

type addInput struct {
	Notifications *notification `json:"notifications" binding:"required"`
	Theme         *string       `json:"theme" binding:"required,oneof=dark light"`
}

func (h handler) AddSettings(c *gin.Context) {
	uid, err := utils.ExtractTokenID(c)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var input addInput

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

	now := time.Now()

	s := models.Settings{
		UserId: uid,
		Notifications: &models.Notification{
			Email:  input.Notifications.Email,
			Mobile: input.Notifications.Mobile,
		},
		Theme:     input.Theme,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	settingsCollection := h.DB.Collection("settings")
	req, err := settingsCollection.InsertOne(context.TODO(), s)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "settings added successfully",
		"id":      req.InsertedID,
	})
}
