package tasks

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) GetTask(c *gin.Context) {
	c.String(http.StatusOK, "get task route")
}
