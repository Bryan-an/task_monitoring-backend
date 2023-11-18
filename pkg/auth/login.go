package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Bryan-an/tasker-backend/pkg/common/models"
	"github.com/Bryan-an/tasker-backend/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"golang.org/x/crypto/bcrypt"
)

type loginInput struct {
	Email    *string `json:"email" binding:"required,email"`
	Password *string `json:"password" binding:"required"`
}

func (h handler) Login(c *gin.Context) {
	var input loginInput

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

	usersCollection := h.DB.Collection("users")

	filter := bson.D{
		{Key: "email", Value: input.Email},
		{Key: "status", Value: "active"},
	}

	var u models.User

	if err := usersCollection.FindOne(context.TODO(), filter).Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "user or password incorrect",
			})

			return
		}

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if correct := verifyPassword(*input.Password, *u.Password); !correct {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "user or password incorrect",
		})

		return
	}

	token, err := utils.GenerateToken(u.Id.Hex())

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func SignInUser(details models.UserDetails, db *mongo.Database, client *mongo.Client) (string, error) {
	if details == (models.UserDetails{}) {
		return "", errors.New("user details can't be empty")
	}

	if details.Email == "" {
		return "", errors.New("email can't be empty")
	}

	if details.Name == "" {
		details.Name = details.Email
	}

	usersCollection := db.Collection("users")

	filter := bson.D{
		{Key: "email", Value: details.Email},
		{Key: "status", Value: "active"},
	}

	var user models.User
	var token string
	var tokenErr error

	if err := usersCollection.FindOne(context.TODO(), filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			wc := writeconcern.Majority()
			txnOptions := options.Transaction().SetWriteConcern(wc)
			session, err := client.StartSession()

			if err != nil {
				return "", err
			}

			defer session.EndSession(context.TODO())

			_, err = session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
				role := "user"
				status := "active"
				now := time.Now()

				u := models.User{
					Name:      &details.Name,
					Email:     &details.Email,
					Role:      &role,
					Status:    &status,
					CreatedAt: &now,
					UpdatedAt: &now,
				}

				req, err := usersCollection.InsertOne(ctx, u)

				if err != nil {
					return "", errors.New("error occurred while registering user")
				}

				uid := req.InsertedID.(primitive.ObjectID)

				emailNotifications := false
				mobileNotifications := true
				theme := "light"

				s := models.Settings{
					UserId: &uid,
					Notifications: &models.Notification{
						Email:  &emailNotifications,
						Mobile: &mobileNotifications,
					},
					Theme:     &theme,
					CreatedAt: &now,
					UpdatedAt: &now,
				}

				settingsCollection := db.Collection("settings")

				if _, err = settingsCollection.InsertOne(ctx, s); err != nil {
					return "", errors.New("error occurred while registering user")
				}

				token, tokenErr = utils.GenerateToken(uid.Hex())
				return "", nil
			}, txnOptions)

			if err != nil {
				return "", err
			}
		} else {
			return "", errors.New("error occurred while logging in user")
		}

	} else {
		token, tokenErr = utils.GenerateToken(user.Id.Hex())
	}

	if tokenErr != nil {
		return "", errors.New("error occurred while generating auth token")
	}

	return token, nil
}
