package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) ReplaceUser(c *gin.Context) {
	c.String(http.StatusOK, "replace user route")
}
