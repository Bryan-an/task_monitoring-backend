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

	from := time.Now()
	to := time.Now().Add(24 * time.Hour)

	filter := bson.M{
		"user_id": uid,
		"status":  "created",
		"date": bson.M{
			"$gte": primitive.NewDateTimeFromTime(
				time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location()).UTC(),
			),
			"$lt": primitive.NewDateTimeFromTime(
				time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, to.Location()).UTC(),
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
		"data": tasks,
	})
}
