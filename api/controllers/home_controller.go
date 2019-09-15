package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Todo item created successfully!",
	})
}
