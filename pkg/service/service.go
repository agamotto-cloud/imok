package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Healz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"health": "well",
	})
}
