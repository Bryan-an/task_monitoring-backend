package users

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
	Name *string `json:"name"`
}

func (h handler) UpdateUser(c *gin.Context) {
	uid, err := utils.ExtractTokenID(c)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	id, err := primitive.ObjectIDFromHex(uid)

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

	coll := h.DB.Collection("users")

	filter := bson.D{
		{Key: "_id", Value: id},
		{Key: "status", Value: "active"},
	}

	data := bson.M{
		"updated_at": time.Now(),
	}

	if input.Name != nil {
		data["name"] = input.Name
	}

	update := bson.D{
		{
			Key:   "$set",
			Value: data,
		},
	}

	result, err := coll.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	if result.MatchedCount == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("user not found with id '%s'", uid),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messasge": "user info updated successfully",
	})
}
