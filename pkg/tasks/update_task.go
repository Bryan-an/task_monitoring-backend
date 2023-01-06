package tasks

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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type updateInput struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Labels      *[]string  `json:"labels"`
	Priority    *string    `json:"priority"`
	Complexity  *string    `json:"complexity"`
	From        *time.Time `json:"from"`
	To          *time.Time `json:"to"`
	Done        *bool      `json:"done"`
}

func (h handler) UpdateTask(c *gin.Context) {
	taskId := c.Param("id")
	uid, err := utils.ExtractTokenID(c)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	id, err := primitive.ObjectIDFromHex(taskId)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	var input updateInput

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

	tasksCollection := h.DB.Collection("tasks")

	filter := bson.D{
		{Key: "_id", Value: id},
		{Key: "user_id", Value: uid},
		{Key: "status", Value: "created"},
	}

	data := bson.M{
		"updated_at": time.Now(),
	}

	if input.Title != nil {
		data["title"] = input.Title
	}

	if input.Description != nil {
		data["description"] = input.Description
	}

	if input.Labels != nil {
		data["labels"] = input.Labels
	}

	if input.Priority != nil {
		data["priority"] = input.Priority
	}

	if input.Complexity != nil {
		data["complexity"] = input.Complexity
	}

	if input.From != nil {
		data["from"] = input.From
	}

	if input.To != nil {
		data["to"] = input.To
	}

	if input.Done != nil {
		data["done"] = input.Done
	}

	update := bson.D{
		{
			Key:   "$set",
			Value: data,
		},
	}

	result, err := tasksCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	if result.MatchedCount == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("task not found with id '%s'", uid),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "task updated successfully",
	})
}
