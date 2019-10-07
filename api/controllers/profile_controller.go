package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
)

func (server *Server) UpdateProfile(c *gin.Context){
	file, _ := c.FormFile("file")
	log.Println(file.Filename)

}