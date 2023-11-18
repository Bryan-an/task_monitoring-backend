package auth

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	gomail "gopkg.in/mail.v2"

	"github.com/Bryan-an/tasker-backend/pkg/common/models"
	"github.com/Bryan-an/tasker-backend/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type registerInput struct {
	Name     *string `json:"name"`
	Email    *string `json:"email" binding:"required,email"`
	Password *string `json:"password" binding:"required"`
}

func (h handler) Register(c *gin.Context) {
	var input registerInput

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

	err := validateUser(h, c, input)

	if err != nil {
		return
	}

	if input.Name == nil {
		input.Name = input.Email
	}

	hash, err := utils.HashPassword(*input.Password)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	input.Password = &hash

	role := "user"
	status := "created"
	now := time.Now()

	u := models.User{
		Name:      input.Name,
		Email:     input.Email,
		Password:  input.Password,
		Role:      &role,
		Status:    &status,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	usersCollection := h.DB.Collection("users")

	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)
	session, err := h.Client.StartSession()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	defer session.EndSession(context.TODO())

	result, err := session.WithTransaction(
		context.TODO(),
		func(ctx mongo.SessionContext) (interface{}, error) {
			uResult, err := usersCollection.InsertOne(ctx, u)

			if err != nil {
				return uResult, err
			}

			userId := uResult.InsertedID.(primitive.ObjectID).Hex()
			emailNotifications := false
			mobileNotifications := true
			theme := "light"
			now := time.Now()

			s := models.Settings{
				UserId: &userId,
				Notifications: &models.Notification{
					Email:  &emailNotifications,
					Mobile: &mobileNotifications,
				},
				Theme:     &theme,
				CreatedAt: &now,
				UpdatedAt: &now,
			}

			settingsCollection := h.DB.Collection("settings")

			if result, err := settingsCollection.InsertOne(ctx, s); err != nil {
				return result, err
			}

			err = SendVerificationEmail(h, c, u, ctx)

			if err != nil {
				return nil, err
			}

			return uResult, nil
		}, txnOptions)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "please check your email for email verification code",
		"id":      result.(*mongo.InsertOneResult).InsertedID,
	})
}

func validateUser(h handler, c *gin.Context, input registerInput) error {
	usersCollection := h.DB.Collection("users")

	filter := bson.D{
		{Key: "email", Value: input.Email},
		{Key: "status", Value: "active"},
	}

	var user models.User

	if err := usersCollection.FindOne(context.TODO(), filter).Decode(&user); err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "this email address is already in use",
		})

		return errors.New("this email address is already in use")
	}

	minEntropy, err := strconv.ParseFloat(os.Getenv("MIN_ENTROPY_BITS"), 64)

	if err != nil {
		minEntropy = 50
	}

	if err := passwordvalidator.Validate(*input.Password, minEntropy); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []utils.ErrorMsg{
			{
				Field:   "Password",
				Message: err.Error(),
			},
		}})

		return err
	}

	return nil
}

func SendVerificationEmail(h handler, c *gin.Context, user models.User, msc mongo.SessionContext) error {
	m := gomail.NewMessage()
	from := os.Getenv("SENDER_EMAIL")
	to := *user.Email
	password := os.Getenv("SENDER_PASSWORD")
	host := "smtp.gmail.com"
	port := 587
	otp, err := utils.GetOTPToken(6)

	if err != nil {
		return err
	}

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Tasker - Email code verification")
	m.SetBody("text/html", "<p>This is your email verification code for Tasker: <b>"+otp+"</b></p>")
	d := gomail.NewDialer(host, port, from, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	lifespan, err := strconv.Atoi(os.Getenv("EMAIL_VERIFICATION_CODE_EXPIRATION"))

	if err != nil {
		return err
	}

	expiresAt := time.Now().Add(time.Second * time.Duration(lifespan))

	data := &models.VerificationData{
		Email:     user.Email,
		Code:      &otp,
		ExpiresAt: &expiresAt,
	}

	coll := h.DB.Collection("verifications")

	if _, err = coll.InsertOne(msc, data); err != nil {
		return err
	}

	return nil
}
