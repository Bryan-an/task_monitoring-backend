package users

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

func (h handler) DeleteUser(c *gin.Context) {
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

	coll := h.DB.Collection("users")

	filter := bson.D{
		{Key: "_id", Value: id},
		{Key: "status", Value: "active"},
	}

	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "status", Value: "unsubscribed"},
				{Key: "updated_at", Value: time.Now()},
			},
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
		"messasge": "user unsubscribed",
	})
}
