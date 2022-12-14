package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) GetUser(c *gin.Context) {
	c.String(http.StatusOK, "get user route")
}
