package tasks

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) UpdateTask(c *gin.Context) {
	c.String(http.StatusOK, "update task route")
}
