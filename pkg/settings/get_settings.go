package settings

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Bryan-an/tasker-backend/pkg/common/models"
	"github.com/Bryan-an/tasker-backend/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h handler) GetSettings(c *gin.Context) {
	uid, err := utils.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Fatal(err)
		return
	}

	settingsCollection := h.DB.Collection("settings")
	var settings models.Settings
	filter := bson.D{{Key: "user_id", Value: uid}}

	if err = settingsCollection.FindOne(context.TODO(), filter).Decode(&settings); err != nil {
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("settings not found for user with id '%s'", uid),
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Fatal(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"settings": settings})
}
