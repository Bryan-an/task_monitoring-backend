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
	facebookOAuth "golang.org/x/oauth2/facebook"
)

type loginInputMobile struct {
	Token *string `json:"token" binding:"required"`
}

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

func (h handler) LoginWithFacebookMobile(c *gin.Context) {
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

	details, err := GetUserInfoFromFacebook(*input.Token)

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
