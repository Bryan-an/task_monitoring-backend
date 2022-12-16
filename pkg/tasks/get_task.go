package tasks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Bryan-an/tasker-backend/pkg/common/models"
	"github.com/Bryan-an/tasker-backend/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h handler) GetTask(c *gin.Context) {
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
	var task models.Task

	filter := bson.D{
		{Key: "user_id", Value: uid},
		{Key: "_id", Value: id},
		{Key: "status", Value: "created"},
	}

	if err = tasksCollection.FindOne(context.TODO(), filter).Decode(&task); err != nil {
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("task not found with id '%s'", taskId),
			})

			return
		}

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}
