package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/victorsteven/fullstack/api/auth"
	"github.com/victorsteven/fullstack/api/models"
	"github.com/victorsteven/fullstack/api/utils/formaterror"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (server *Server) LikePost(c *gin.Context){

	//clear previous error if any
	errList = map[string]string{}

	// Start processing the request
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	requestBody := make(map[string]interface{})
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  err,
		})
		return
	}
	//convert the interface to string
	formatedRequest := map[string]int{}
	for key, value := range requestBody {
		strKey := fmt.Sprintf("%v", key)
		intValue := int(value.(float64))
		formatedRequest[strKey] = intValue
	}

	uid := uint32(formatedRequest["user_id"])
	pid := uint64(formatedRequest["post_id"])
	like := uint64(formatedRequest["like"])

	// Get user id from the token for valid tokens
	tokenID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// If the id is not the authenticated user id
	if tokenID != 0 && tokenID != uid {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	likeModel := models.Like{}
	likeModel.UserID = uid
	likeModel.PostID = pid
	likeModel.Like = like

	likeCreated, err := likeModel.SaveLike(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		errList = formattedError
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": likeCreated,
	})
}

func (server *Server) UnLikePost(c *gin.Context){

	//clear previous error if any
	errList = map[string]string{}

	// Start processing the request
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	requestBody := make(map[string]interface{})
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  err,
		})
		return
	}
	//convert the interface to string
	formatedRequest := map[string]int{}
	for key, value := range requestBody {
		strKey := fmt.Sprintf("%v", key)
		intValue := int(value.(float64))
		formatedRequest[strKey] = intValue
	}
	uid := uint32(formatedRequest["user_id"])
	pid := uint64(formatedRequest["post_id"])
	like := uint64(formatedRequest["like"])

	// Get user id from the token for valid tokens
	tokenID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// If the id is not the authenticated user id
	if tokenID != 0 && tokenID != uid {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	likeModel := models.Like{}
	likeModel.UserID = uid
	likeModel.PostID = pid
	likeModel.Like = like

	_, err = likeModel.DeleteLike(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		errList = formattedError
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{
		"status":   http.StatusNoContent,
		"response": "Deleted like",
	})
}

func (server *Server) GetLikes(c *gin.Context){
	postID := c.Param("id")

	//fmt.Println("THe post is is id: ", postID)
	// Is a valid post id given to us?
	pid, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	like := models.Like{}

	likes := like.GetLikesCount(server.DB, pid)

	fmt.Printf("this is the likes received: %d\n", likes)

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": likes,
	})
}