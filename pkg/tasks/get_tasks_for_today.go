package tasks

import (
	"context"
	"net/http"
	"time"

	"github.com/Bryan-an/tasker-backend/pkg/common/models"
	"github.com/Bryan-an/tasker-backend/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h handler) GetTasksForToday(c *gin.Context) {
	uid, err := utils.ExtractTokenID(c)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tasksCollection := h.DB.Collection("tasks")
	var tasks []models.Task

	filter := bson.M{
		"user_id": uid,
		"status":  "created",
		"date": bson.M{
			"$gte": primitive.NewDateTimeFromTime(
				time.Now().UTC().Truncate(24 * time.Hour),
			),
			"$lt": primitive.NewDateTimeFromTime(
				time.Now().UTC().Add(24 * time.Hour).Truncate(24 * time.Hour),
			),
		},
	}

	opts := options.Find().SetSort(bson.D{{Key: "updated_at", Value: -1}})

	cursor, err := tasksCollection.Find(context.TODO(), filter, opts)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err = cursor.All(context.TODO(), &tasks); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
	})
}
