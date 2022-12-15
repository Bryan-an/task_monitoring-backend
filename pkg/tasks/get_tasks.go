package tasks

import (
	"context"
	"net/http"

	"github.com/Bryan-an/tasker-backend/pkg/common/models"
	"github.com/Bryan-an/tasker-backend/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) GetTasks(c *gin.Context) {
	uid, err := utils.ExtractTokenID(c)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	taskCollection := h.DB.Collection("tasks")
	var tasks []models.Task

	filter := bson.D{
		{Key: "user_id", Value: uid},
		{Key: "status", Value: "created"},
	}

	cursor, err := taskCollection.Find(context.TODO(), filter)

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

	c.JSON(http.StatusOK, gin.H{"tasks": tasks, "count": len(tasks)})
}
