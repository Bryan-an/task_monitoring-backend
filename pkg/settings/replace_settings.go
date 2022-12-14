package settings

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) ReplaceSettings(c *gin.Context) {
	c.String(http.StatusOK, "replace settings route")
}
