package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Errors.JSON() != nil {
			for _, ginErr := range c.Errors {
				log.Println(ginErr.Err.Error())
			}

			c.JSON(-1, c.Errors.JSON())
		}
	}
}
