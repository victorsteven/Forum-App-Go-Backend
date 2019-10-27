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



func (server *Server) LikePost(c *gin.Context) {
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
	_, err = auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	//convert the interface type to integer
	userID  := uint32((requestBody["user_id"]).(float64))
	//postID  := int((requestBody["post_id"]).(float64))
	postString := fmt.Sprintf("%v", requestBody["post_id"]) //convert interface to string
	postInt, _ := strconv.Atoi(postString) //convert string to integer
	postID := uint64(postInt)

	like := models.Like{}
	like.UserID = userID
	like.PostID = postID

	likeCreated, err := like.SaveLike(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		errList = formattedError
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	//data := make(map[string]interface{})
	//data["postID"] = postID
	//data["likes"] = likeCreated

	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": likeCreated,
	})
}

func (server *Server) UnLikePost(c *gin.Context) {

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
	tokenID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	//convert the interface type to integer
	likeID  := uint64((requestBody["id"]).(float64))
	//postID  := int((requestBody["post_id"]).(float64))
	userID  := uint32((requestBody["user_id"]).(float64))
	postString := fmt.Sprintf("%v", requestBody["post_id"]) //convert interface to string
	postInt, _ := strconv.Atoi(postString) //convert string to integer
	postID := uint64(postInt)

	//likeString := fmt.Sprintf("%v", requestBody["like_id"]) //convert interface to string
	//likeInt, _ := strconv.Atoi(likeString) //convert string to integer
	//likeID := uint64(likeInt)

	// If the id is not the authenticated user id
	if tokenID != 0 && tokenID != userID {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	like := models.Like{}
	like.ID = likeID
	like.UserID = userID
	like.PostID = postID

	likeDeleted, err := like.DeleteLike(server.DB, likeID)
	if err != nil {
		//formattedError := formaterror.FormatError(err.Error())
		//errList = formattedError
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": likeDeleted,
	})
}

func (server *Server) GetLikes(c *gin.Context){

	//clear previous error if any
	errList = map[string]string{}

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

	likes, err := like.GetLikesInfo(server.DB, pid)
	if err != nil {
		errList["No_likes"] = "No Likes found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	//fmt.Printf("this is the likes received: %d\n", len(likes))

	//data := make(map[string]interface{})
	//data["postID"] = pid
	//data["likes"] = likes

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": likes,
	})
}

//func (server *Server) getAuthUserLike(c *gin.Context){
//
//	postID := c.Param("id")
//	// Is a valid post id given to us?
//	pid, err := strconv.ParseUint(postID, 10, 64)
//	if err != nil {
//		errList["Invalid_request"] = "Invalid Request"
//		c.JSON(http.StatusBadRequest, gin.H{
//			"status": http.StatusBadRequest,
//			"error":  errList,
//		})
//		return
//	}
//
//	// Is this user authenticated?
//	uid, err := auth.ExtractTokenID(c.Request)
//	if err != nil {
//		errList["Unauthorized"] = "Unauthorized"
//		c.JSON(http.StatusUnauthorized, gin.H{
//			"status": http.StatusUnauthorized,
//			"error":  errList,
//		})
//		return
//	}
//}