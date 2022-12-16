package tasks

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

type AddInput struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Labels      []string `json:"labels"`
	Priority    string   `json:"priority" binding:"required,oneof=low medium high"`
	Complexity  string   `json:"complexity" binding:"required,oneof=low medium high"`
}

func (h handler) AddTask(c *gin.Context) {
	uid, err := utils.ExtractTokenID(c)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	var input AddInput

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

	t := models.Task{
		UserId:      uid,
		Title:       input.Title,
		Description: input.Description,
		Labels:      input.Labels,
		Priority:    input.Priority,
		Complexity:  input.Complexity,
		Status:      "created",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tasksCollection := h.DB.Collection("tasks")
	req, err := tasksCollection.InsertOne(context.TODO(), t)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "task added successfully",
		"id":      req.InsertedID,
	})
}
