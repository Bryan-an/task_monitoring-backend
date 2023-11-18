package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/Bryan-an/tasker-backend/pkg/common/models"
	"github.com/Bryan-an/tasker-backend/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GetGoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Endpoint:     google.Endpoint,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}
}

func GetUserInfoFromGoogle(token string) (models.UserDetails, error) {
	var details models.UserDetails

	req, _ := http.NewRequest(
		"GET",
		"https://www.googleapis.com/oauth2/v2/userinfo?access_token="+token,
		nil,
	)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return models.UserDetails{},
			errors.New("error ocurred while getting user info from Google")
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&details)
	defer res.Body.Close()

	if err != nil {
		return models.UserDetails{},
			errors.New("error ocurred while getting user info from Google")
	}

	return details, nil
}

func (h handler) InitGoogleLogin(c *gin.Context) {
	var config = GetGoogleOAuthConfig()
	url := config.AuthCodeURL(GetRandomOAuthStateString())
	http.Redirect(c.Writer, c.Request, url, http.StatusTemporaryRedirect)
}

func (h handler) HandleGoogleLogin(c *gin.Context) {
	var state = c.Request.FormValue("state")
	var code = c.Request.FormValue("code")

	if state != GetRandomOAuthStateString() {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": "error while logging in user"},
		)

		return
	}

	if code == "" {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": "error while logging in user"},
		)

		return
	}

	var config = GetGoogleOAuthConfig()
	token, err := config.Exchange(context.TODO(), code)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	details, err := GetUserInfoFromGoogle(token.AccessToken)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	authToken, err := SignInUser(details, h.DB, h.Client)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"token": authToken})
}

func (h handler) LoginWithGoogleMobile(c *gin.Context) {
	var input loginInputMobile

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

	details, err := GetUserInfoFromGoogle(*input.Token)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	authToken, err := SignInUser(details, h.DB, h.Client)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"token": authToken})
}
