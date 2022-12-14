package tasks

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) DeleteTask(c *gin.Context) {
	c.String(http.StatusOK, "delete task route")
}
