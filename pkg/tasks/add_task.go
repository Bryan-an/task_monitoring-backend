package tasks

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) AddTask(c *gin.Context) {
	c.String(http.StatusOK, "add task route")
}
