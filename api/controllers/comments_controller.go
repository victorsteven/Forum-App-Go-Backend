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

//func (server *Server) CreateCommente(c *gin.Context){
//
//	//clear previous error if any
//	errList = map[string]string{}
//
//	body, err := ioutil.ReadAll(c.Request.Body)
//	if err != nil {
//		errList["Invalid_body"] = "Unable to get request"
//		c.JSON(http.StatusUnprocessableEntity, gin.H{
//			"status": http.StatusUnprocessableEntity,
//			"error":  errList,
//		})
//		return
//	}
//
//	comment := models.Comment{}
//
//	err = json.Unmarshal(body, &comment)
//	if err != nil {
//		errList["Unmarshal_error"] = "Cannot unmarshal body"
//		c.JSON(http.StatusUnprocessableEntity, gin.H{
//			"status": http.StatusUnprocessableEntity,
//			"error":  errList,
//		})
//		return
//	}
//
//	comment.Prepare()
//	errorMessages := comment.Validate("")
//	if len(errorMessages) > 0 {
//		errList = errorMessages
//		c.JSON(http.StatusUnprocessableEntity, gin.H{
//			"status": http.StatusUnprocessableEntity,
//			"error":  errList,
//		})
//		return
//	}
//
//	uid, err := auth.ExtractTokenID(c.Request)
//	if err != nil {
//		errList["Unauthorized"] = "Unauthorized"
//		c.JSON(http.StatusUnauthorized, gin.H{
//			"status": http.StatusUnauthorized,
//			"error":  errList,
//		})
//		return
//	}
//
//	if uid != comment.UserID {
//		errList["Unauthorized"] = "Unauthorized"
//		c.JSON(http.StatusUnauthorized, gin.H{
//			"status": http.StatusUnauthorized,
//			"error":  errList,
//		})
//		return
//	}
//
//	commetCreated, err := comment.SaveComment(server.DB)
//	if err != nil {
//		formattedError := formaterror.FormatError(err.Error())
//		errList = formattedError
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"status": http.StatusInternalServerError,
//			"error":  errList,
//		})
//		return
//	}
//	c.JSON(http.StatusCreated, gin.H{
//		"status":   http.StatusCreated,
//		"response": commetCreated,
//	})
//}


func (server *Server) CreateComment(c *gin.Context) {
	//clear previous error if any
	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	//the value of the post id and the user id are integers
	requestBody := make(map[string]interface{})
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	comment := models.Comment{}

	//convert the interface type to integer
	userID  := uint32((requestBody["user_id"]).(float64))
	postString := fmt.Sprintf("%v", requestBody["post_id"]) //convert interface to string
	postInt, _ := strconv.Atoi(postString) //convert string to integer
	postID := uint64(postInt)
	commentString := fmt.Sprintf("%v", requestBody["body"])

	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	if uid != userID {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	comment.UserID = userID
	comment.PostID = postID
	comment.Body = commentString

	comment.Prepare()
	errorMessages := comment.Validate("")
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	commentCreated, err := comment.SaveComment(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		errList = formattedError
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": commentCreated,
	})
}

func (server *Server) GetComments(c *gin.Context){

	//clear previous error if any
	errList = map[string]string{}

	postID := c.Param("id")

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
	comment := models.Comment{}

	comments, err := comment.GetComments(server.DB, pid)
	if err != nil {
		errList["No_comments"] = "No Likes found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": comments,
	})
}