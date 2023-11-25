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
	Title       utils.JSONString      `json:"title"`
	Description utils.JSONString      `json:"description"`
	Labels      utils.JSONStringSlice `json:"labels"`
	Priority    utils.JSONString      `json:"priority"`
	Complexity  utils.JSONString      `json:"complexity"`
	Date        utils.JSONTime        `json:"date"`
	From        utils.JSONTime        `json:"from"`
	To          utils.JSONTime        `json:"to"`
	Done        utils.JSONBool        `json:"done"`
	Remind      utils.JSONBool        `json:"remind"`
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

	if input.Title.Set {
		if input.Title.Valid {
			data["title"] = input.Title.Value
		} else {
			data["title"] = nil
		}
	}

	if input.Description.Set {
		if input.Description.Valid {
			data["description"] = input.Description.Value
		} else {
			data["description"] = nil
		}
	}

	if input.Labels.Set {
		if input.Labels.Valid {
			data["labels"] = input.Labels.Value
		} else {
			data["labels"] = nil
		}
	}

	if input.Priority.Set {
		if input.Priority.Valid {
			data["priority"] = input.Priority.Value
		} else {
			data["priority"] = nil
		}
	}

	if input.Complexity.Set {
		if input.Complexity.Valid {
			data["complexity"] = input.Complexity.Value
		} else {
			data["complexity"] = nil
		}
	}

	if input.Date.Set {
		if input.Date.Valid {
			data["date"] = input.Date.Value
		} else {
			data["date"] = nil
		}
	}

	if input.From.Set {
		if input.From.Valid {
			data["from"] = input.From.Value
		} else {
			data["from"] = nil
		}
	}

	if input.To.Set {
		if input.To.Valid {
			data["to"] = input.To.Value
		} else {
			data["to"] = nil
		}
	}

	if input.Done.Set {
		if input.Done.Valid {
			data["done"] = input.Done.Value
		} else {
			data["done"] = nil
		}
	}

	if input.Remind.Set {
		if input.Remind.Valid {
			data["remind"] = input.Remind.Value
		} else {
			data["remind"] = nil
		}
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
