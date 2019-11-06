package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/victorsteven/forum/api/auth"
	"github.com/victorsteven/forum/api/models"
	"github.com/victorsteven/forum/api/utils/formaterror"
)

func (server *Server) CreatePost(c *gin.Context) {

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
	post := models.Post{}

	err = json.Unmarshal(body, &post)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// check if the user exist:
	user := models.User{}
	err = server.DB.Debug().Model(models.User{}).Where("id = ?", uid).Take(&user).Error
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	post.AuthorID = uid //the authenticated user is the one creating the post

	post.Prepare()
	errorMessages := post.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	postCreated, err := post.SavePost(server.DB)
	if err != nil {
		errList := formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": postCreated,
	})
}

func (server *Server) GetPosts(c *gin.Context) {

	post := models.Post{}

	posts, err := post.FindAllPosts(server.DB)
	if err != nil {
		errList["No_post"] = "No Post Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": posts,
	})
}

func (server *Server) GetPost(c *gin.Context) {

	postID := c.Param("id")
	pid, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	post := models.Post{}

	postReceived, err := post.FindPostByID(server.DB, pid)
	if err != nil {
		errList["No_post"] = "No Post Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": postReceived,
	})
}

func (server *Server) UpdatePost(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	postID := c.Param("id")
	// Check if the post id is valid
	pid, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	//Check if the post exist
	origPost := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&origPost).Error
	if err != nil {
		errList["No_post"] = "No Post Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	if uid != origPost.AuthorID {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// Read the data posted
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	// Start processing the request data
	post := models.Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	post.ID = origPost.ID //this is important to tell the model the post id to update, the other update field are set above
	post.AuthorID = origPost.AuthorID

	post.Prepare()
	errorMessages := post.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	postUpdated, err := post.UpdateAPost(server.DB)
	if err != nil {
		errList := formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": postUpdated,
	})
}

func (server *Server) DeletePost(c *gin.Context) {

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

	fmt.Println("this is delete post sir")

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// Check if the post exist
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		errList["No_post"] = "No Post Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	// Is the authenticated user, the owner of this post?
	if uid != post.AuthorID {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// If all the conditions are met, delete the post
	_, err = post.DeleteAPost(server.DB)
	if err != nil {
		errList["Other_error"] = "Please try again later"
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	comment := models.Comment{}
	like := models.Like{}

	// Also delete the likes and the comments that this post have:
	_, err = comment.DeletePostComments(server.DB, pid)
	if err != nil {
		errList["Other_error"] = "Please try again later"
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	_, err = like.DeletePostLikes(server.DB, pid)
	if err != nil {
		errList["Other_error"] = "Please try again later"
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "Post deleted",
	})
}

func (server *Server) GetUserPosts(c *gin.Context) {

	userID := c.Param("id")
	// Is a valid user id given to us?
	uid, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	post := models.Post{}
	posts, err := post.FindUserPosts(server.DB, uint32(uid))
	if err != nil {
		errList["No_post"] = "No Post Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": posts,
	})
}
