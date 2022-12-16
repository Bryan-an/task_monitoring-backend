package tasks

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Bryan-an/tasker-backend/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h handler) DeleteTask(c *gin.Context) {
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

	tasksCollection := h.DB.Collection("tasks")

	filter := bson.D{
		{Key: "user_id", Value: uid},
		{Key: "_id", Value: id},
		{Key: "status", Value: "created"},
	}

	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "status", Value: "deleted"},
				{Key: "updated_at", Value: time.Now()},
			},
		},
	}

	result, err := tasksCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	if result.MatchedCount == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("task not found with id '%s'", taskId),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "task deleted successfully",
	})
}
