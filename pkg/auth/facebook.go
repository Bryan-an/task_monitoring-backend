package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/Bryan-an/tasker-backend/pkg/common/models"
	"github.com/Bryan-an/tasker-backend/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	facebookOAuth "golang.org/x/oauth2/facebook"
)

func GetFacebookOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
		ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("FACEBOOK_REDIRECT_URL"),
		Endpoint:     facebookOAuth.Endpoint,
		Scopes:       []string{"email"},
	}
}

func GetRandomOAuthStateString() string {
	return os.Getenv("OAUTH_STATE_STRING")
}

func GetUserInfoFromFacebook(token string) (models.UserDetails, error) {
	var details models.UserDetails
	req, _ := http.NewRequest(
		"GET",
		"https://graph.facebook.com/me?fields=id,name,email&access_token="+token,
		nil,
	)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return models.UserDetails{},
			errors.New("error ocurred while getting user info from Facebook")
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&details)
	defer res.Body.Close()

	if err != nil {
		return models.UserDetails{},
			errors.New("error ocurred while getting user info from Facebook")
	}

	return details, nil
}

func (h handler) InitFacebookLogin(c *gin.Context) {
	var config = GetFacebookOAuthConfig()
	url := config.AuthCodeURL(GetRandomOAuthStateString())
	http.Redirect(c.Writer, c.Request, url, http.StatusTemporaryRedirect)
}

func (h handler) HandleFacebookLogin(c *gin.Context) {
	var state = c.Request.FormValue("state")
	var code = c.Request.FormValue("code")

	if state != GetRandomOAuthStateString() {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": "error while logging in user"},
		)

		return
	}

	var config = GetFacebookOAuthConfig()
	token, err := config.Exchange(context.TODO(), code)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	details, err := GetUserInfoFromFacebook(token.AccessToken)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	authToken, err := SignInUser(details, h.DB)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"token": authToken})
}

func SignInUser(details models.UserDetails, db *mongo.Database) (string, error) {
	if details == (models.UserDetails{}) {
		return "", errors.New("user details can't be empty")
	}

	if details.Email == "" {
		return "", errors.New("email can't be empty")
	}

	if details.Name == "" {
		return "", errors.New("name can't be empty'")
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
			u := models.User{
				Name:      details.Name,
				Email:     details.Email,
				Role:      "user",
				Status:    "active",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			req, err := usersCollection.InsertOne(context.TODO(), u)

			if err != nil {
				return "", errors.New("error occurred while registering user")
			}

			uid := req.InsertedID.(primitive.ObjectID).Hex()

			s := models.Settings{
				UserId: uid,
				Notifications: models.Notification{
					Email:  false,
					Mobile: true,
				},
				Security:  "something",
				Theme:     "light",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			settingsCollection := db.Collection("settings")

			if _, err = settingsCollection.InsertOne(context.TODO(), s); err != nil {
				return "", errors.New("error occurred while registering user")
			}

			token, tokenErr = utils.GenerateToken(uid)
		} else {
			return "", errors.New("error occurred while registering user")
		}

	} else {
		token, tokenErr = utils.GenerateToken(user.Id.Hex())
	}

	if tokenErr != nil {
		return "", errors.New("error occurred while generating auth token")
	}

	return token, nil
}
