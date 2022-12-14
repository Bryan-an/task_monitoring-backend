package settings

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) UpdateSettings(c *gin.Context) {
	c.String(http.StatusOK, "update settings route")
}
