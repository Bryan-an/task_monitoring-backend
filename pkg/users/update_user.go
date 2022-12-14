package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) UpdateUser(c *gin.Context) {
	c.String(http.StatusOK, "update user route")
}
