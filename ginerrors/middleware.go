package ginerrors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errors": c.Errors,
			})
			c.Abort()
		}
	}
}
