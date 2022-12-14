package tasks

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) ReplaceTask(c *gin.Context) {
	c.String(http.StatusOK, "replace task route")
}
