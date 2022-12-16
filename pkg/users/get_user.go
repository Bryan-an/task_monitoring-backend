package users

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
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h handler) GetUser(c *gin.Context) {
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
	var user models.User
	opts := options.FindOne().SetProjection(bson.D{{Key: "password", Value: 0}})

	filter := bson.D{
		{Key: "_id", Value: id},
		{Key: "status", Value: "active"},
	}

	if err = coll.FindOne(context.TODO(), filter, opts).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("user not found with id '%s'", uid),
			})

			return
		}

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
