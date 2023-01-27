package tasks

import (
	"bytes"
	"context"
	"net/http"
	"strings"

	"github.com/Bryan-an/tasker-backend/pkg/common/models"
	"github.com/Bryan-an/tasker-backend/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h handler) GetTasks(c *gin.Context) {
	priority := c.Query("priority")
	complexity := c.Query("complexity")
	labels := c.Query("labels")
	done := c.Query("done")
	remind := c.Query("remind")
	order := c.DefaultQuery("order", "des")

	uid, err := utils.ExtractTokenID(c)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	taskCollection := h.DB.Collection("tasks")
	var tasks []models.Task

	filter := bson.M{
		"user_id": uid,
		"status":  "created",
	}

	if priority != "" {
		filter["priority"] = priority
	}

	if complexity != "" {
		filter["complexity"] = complexity
	}

	if labels != "" {
		ls := strings.Split(labels, ",")

		var b bytes.Buffer

		for i, l := range ls {
			if i == 0 {
				b.WriteString("(^" + l + "$)")
			} else {
				b.WriteString("|(^" + l + "$)")
			}
		}

		filter["labels"] = bson.D{{
			Key: "$regex", Value: primitive.Regex{Pattern: b.String(), Options: "i"},
		}}
	}

	if done != "" {
		if done == "true" {
			filter["done"] = true
		} else if done == "false" {
			filter["done"] = false
		}
	}

	if remind != "" {
		if remind == "true" {
			filter["remind"] = true
		} else if remind == "false" {
			filter["remind"] = false
		}
	}

	var sort int

	if order == "asc" {
		sort = 1
	} else {
		sort = -1
	}

	opts := options.Find().SetSort(bson.D{{Key: "updated_at", Value: sort}})

	cursor, err := taskCollection.Find(context.TODO(), filter, opts)

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
