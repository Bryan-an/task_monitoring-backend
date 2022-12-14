package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) DeleteUser(c *gin.Context) {
	c.String(http.StatusOK, "delete user route")
}
