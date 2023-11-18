package utils

import (
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GenerateToken(userId string) (string, error) {
	lifespan, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(lifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

func TokenValid(c *gin.Context) error {
	tokenString := ExtractToken(c)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("API_SECRET")), nil
	})

	if err != nil {
		return err
	}

	return nil
}

func ExtractToken(c *gin.Context) string {
	token := c.Request.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}

	return ""
}

func ExtractTokenID(c *gin.Context) (*primitive.ObjectID, error) {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("API_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		uid := fmt.Sprintf("%v", claims["user_id"])

		objId, err := primitive.ObjectIDFromHex(uid)

		if err != nil {
			return nil, err
		}

		return &objId, nil
	}

	return nil, errors.New("invalid token")
}

func GetOTPToken(length int) (string, error) {
	const otpChars = "1234567890"
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)

	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)

	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}
