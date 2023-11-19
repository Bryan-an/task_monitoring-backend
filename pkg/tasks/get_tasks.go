package tasks

import (
	"bytes"
	"context"
	"math"
	"net/http"
	"strconv"
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
	pageParam := c.Query("page")
	pageSizeParam := c.Query("page_size")

	queryParamsErrors := []utils.ErrorMsg{}

	if pageParam == "" {
		queryParamsErrors = append(queryParamsErrors, utils.ErrorMsg{
			Field:   "page",
			Message: "this query param is required",
		})
	}

	page, err := strconv.Atoi(pageParam)

	if err != nil {
		queryParamsErrors = append(queryParamsErrors, utils.ErrorMsg{
			Field:   "page",
			Message: "this query param must be a number",
		})
	}

	if page < 1 {
		queryParamsErrors = append(queryParamsErrors, utils.ErrorMsg{
			Field:   "page",
			Message: "this query param must be greater than 0",
		})
	}

	if pageSizeParam == "" {
		queryParamsErrors = append(queryParamsErrors, utils.ErrorMsg{
			Field:   "page_size",
			Message: "this query param is required",
		})
	}

	pageSize, err := strconv.Atoi(pageSizeParam)

	if err != nil {
		queryParamsErrors = append(queryParamsErrors, utils.ErrorMsg{
			Field:   "page_size",
			Message: "this query param must be a number",
		})
	}

	if pageSize < 1 {
		queryParamsErrors = append(queryParamsErrors, utils.ErrorMsg{
			Field:   "page_size",
			Message: "this query param must be greater than 0",
		})
	}

	if len(queryParamsErrors) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": queryParamsErrors})
		return
	}

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

	opts := options.Find().
		SetSort(bson.D{{Key: "updated_at", Value: sort}}).
		SetLimit(int64(pageSize)).
		SetSkip(int64((page - 1) * pageSize))

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

	totalRecords, err := taskCollection.CountDocuments(context.TODO(), filter)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	var nextPage *int
	var prevPage *int

	if page < totalPages {
		p := page + 1
		nextPage = &p
	} else {
		nextPage = nil
	}

	if page > 1 && page <= totalPages {
		p := page - 1
		prevPage = &p
	} else {
		prevPage = nil
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
		"pagination": gin.H{
			"count":         len(tasks),
			"page":          page,
			"page_size":     pageSize,
			"total_records": totalRecords,
			"total_pages":   totalPages,
			"next_page":     nextPage,
			"prev_page":     prevPage,
		},
	})
}
