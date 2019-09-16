package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victorsteven/fullstack/api/auth"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.TokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"error":  "Unauthorized",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
