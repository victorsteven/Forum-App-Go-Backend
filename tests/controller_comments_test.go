package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCommentPost(t *testing.T) {

	var firstUserEmail, secondUserEmail string
	var firstPostID uint64

	err := refreshUserPostAndCommentTable()
	if err != nil {
		log.Fatal(err)
	}
	users, posts, err := seedUsersAndPosts()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	for _, user := range users {
		if user.ID == 1 {
			firstUserEmail = user.Email
		}
		if user.ID == 2 {
			secondUserEmail = user.Email
		}
	}
	// Get only the first post, which belongs to first user
	for _, post := range posts {
		if post.ID == 2 {
			continue
		}
		firstPostID = post.ID
	}
	// Login both users
	// user 1 and user 2 password are the same, you can change if you want (Note by the time they are hashed and saved in the db, they are different)
	// Note: the value of the user password before it was hashed is "password". so:
	password := "password"

	// Login First User
	tokenInterface1, err := server.SignIn(firstUserEmail, password)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	token1 := tokenInterface1["token"] //get only the token
	firstUserToken := fmt.Sprintf("Bearer %v", token1)

	// Login Second User
	tokenInterface2, err := server.SignIn(secondUserEmail, password)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	token2 := tokenInterface2["token"] //get only the token
	secondUserToken := fmt.Sprintf("Bearer %v", token2)
	fmt.Println("this is the second user token: ", secondUserToken)

	samples := []struct {
		postIDString string
		inputJSON    string
		statusCode   int
		userID       uint32
		postID       uint64
		Body         string
		tokenGiven   string
	}{
		{
			// User 1 can comment on his post
			postIDString: strconv.Itoa(int(firstPostID)), //we need the id as a string
			inputJSON:    `{"body": "comment from user 1"}`,
			statusCode:   201,
			userID:       1,
			postID:       firstPostID,
			Body:         "comment from user 1",
			tokenGiven:   firstUserToken,
		},
		{
			// User 2 can also comment on user 1 post
			postIDString: strconv.Itoa(int(firstPostID)),
			inputJSON:    `{"body":"comment from user 2"}`,
			statusCode:   201,
			userID:       2,
			postID:       firstPostID,
			Body:         "comment from user 2",
			tokenGiven:   secondUserToken,
		},
		{
			// When no body is provided:
			postIDString: strconv.Itoa(int(firstPostID)),
			inputJSON:    `{"body":""}`,
			statusCode:   422,
			postID:       firstPostID,
			tokenGiven:   secondUserToken,
		},
		{
			// Not authenticated (No token provided)
			postIDString: strconv.Itoa(int(firstPostID)),
			statusCode:   401,
			tokenGiven:   "",
		},
		{
			// Wrong Token
			postIDString: strconv.Itoa(int(firstPostID)),
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
		},
		{
			// When invalid post id is given
			postIDString: "unknwon",
			statusCode:   400,
		},
	}

	for _, v := range samples {

		gin.SetMode(gin.TestMode)

		r := gin.Default()

		r.POST("/comments/:id", server.CreateComment)
		req, err := http.NewRequest(http.MethodPost, "/comments/"+v.postIDString, bytes.NewBufferString(v.inputJSON))
		req.Header.Set("Authorization", v.tokenGiven)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 201 {
			responseMap := responseInterface["response"].(map[string]interface{})
			assert.Equal(t, responseMap["post_id"], float64(v.postID))
			assert.Equal(t, responseMap["user_id"], float64(v.userID))
			assert.Equal(t, responseMap["body"], v.Body)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 {
			responseMap := responseInterface["error"].(map[string]interface{})
			if responseMap["Invalid_request"] != nil {
				assert.Equal(t, responseMap["Invalid_request"], "Invalid Request")
			}
			if responseMap["Unauthorized"] != nil {
				assert.Equal(t, responseMap["Unauthorized"], "Unauthorized")
			}
			if responseMap["Required_body"] != nil {
				assert.Equal(t, responseMap["Required_body"], "Required Comment")
			}
		}
	}
}

func TestGetComments(t *testing.T) {
	gin.SetMode(gin.TestMode)
	err := refreshUserPostAndCommentTable()
	if err != nil {
		log.Fatal(err)
	}
	post, users, comments, err := seedUsersPostsAndComments()
	if err != nil {
		log.Fatalf("Cannot seed tables %v\n", err)
	}
	commentsSample := []struct {
		postID         string
		usersLength    int
		commentsLength int
		statusCode     int
	}{
		{
			postID:         strconv.Itoa(int(post.ID)),
			statusCode:     200,
			usersLength:    len(users),
			commentsLength: len(comments),
		},
		{
			postID:     "unknwon",
			statusCode: 400,
		},
		{
			postID:     strconv.Itoa(12322), //an id that does not exist
			statusCode: 404,
		},
	}
	for _, v := range commentsSample {
		r := gin.Default()
		r.GET("/comments/:id", server.GetComments)
		req, err := http.NewRequest(http.MethodGet, "/comments/"+v.postID, nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
		if err != nil {
			t.Errorf("Cannot convert to json here: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			responseMap := responseInterface["response"].([]interface{})
			assert.Equal(t, len(responseMap), v.commentsLength)
			assert.Equal(t, v.usersLength, 2)
		}
		if v.statusCode == 400 || v.statusCode == 404 {
			responseMap := responseInterface["error"].(map[string]interface{})
			if responseMap["Invalid_request"] != nil {
				assert.Equal(t, responseMap["Invalid_request"], "Invalid Request")
			}
			if responseMap["No_post"] != nil {
				assert.Equal(t, responseMap["No_post"], "No post found")
			}
		}
	}
}

func TestUpdateComment(t *testing.T) {

	gin.SetMode(gin.TestMode)

	var secondUserEmail, secondUserPassword string
	var secondUserID uint32
	var secondCommentID uint64

	err := refreshUserPostAndCommentTable()
	if err != nil {
		log.Fatal(err)
	}
	post, users, comments, err := seedUsersPostsAndComments()
	if err != nil {
		log.Fatalf("Cannot seed tables %v\n", err)
	}
	// Get only the second user
	for _, user := range users {
		if user.ID == 1 {
			continue
		}
		secondUserID = user.ID
		secondUserEmail = user.Email
		secondUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	// Get only the second comment
	for _, comment := range comments {
		if comment.ID == 1 {
			continue
		}
		secondCommentID = comment.ID
	}
	//Login the user and get the authentication token
	tokenInterface, err := server.SignIn(secondUserEmail, secondUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	token := tokenInterface["token"] //get only the token
	tokenString := fmt.Sprintf("Bearer %v", token)

	commentsSample := []struct {
		commentID  string
		updateJSON string
		Body       string
		tokenGiven string
		statusCode int
	}{
		{
			commentID:  strconv.Itoa(int(secondCommentID)),
			updateJSON: `{"Body":"This is the update body"}`,
			statusCode: 200,
			Body:       "This is the update body",
			tokenGiven: tokenString,
		},
		{
			// When the body field is empty
			commentID:  strconv.Itoa(int(secondCommentID)),
			updateJSON: `{"Body":""}`,
			statusCode: 422,
			tokenGiven: tokenString,
		},
		{
			//an id that does not exist
			commentID:  strconv.Itoa(12322),
			statusCode: 404,
			tokenGiven: tokenString,
		},
		{
			//When the user is not authenticated
			commentID:  strconv.Itoa(int(secondCommentID)),
			statusCode: 401,
			tokenGiven: "",
		},
		{
			//When wrong token is passed
			commentID:  strconv.Itoa(int(secondCommentID)),
			statusCode: 401,
			tokenGiven: "this is a wrong token",
		},
		{
			// When id passed is invalid
			commentID:  "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range commentsSample {
		r := gin.Default()
		r.PUT("/comments/:id", server.UpdateComment)
		req, err := http.NewRequest(http.MethodPut, "/comments/"+v.commentID, bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		req.Header.Set("Authorization", v.tokenGiven)
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
		if err != nil {
			t.Errorf("Cannot convert to json here: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			responseMap := responseInterface["response"].(map[string]interface{})
			assert.Equal(t, responseMap["post_id"], float64(post.ID))
			assert.Equal(t, responseMap["user_id"], float64(secondUserID))
			assert.Equal(t, responseMap["body"], v.Body)
		}
		if v.statusCode == 400 || v.statusCode == 401 || v.statusCode == 404 {
			responseMap := responseInterface["error"].(map[string]interface{})
			if responseMap["Invalid_request"] != nil {
				assert.Equal(t, responseMap["Invalid_request"], "Invalid Request")
			}
			if responseMap["Unauthorized"] != nil {
				assert.Equal(t, responseMap["Unauthorized"], "Unauthorized")
			}
			if responseMap["No_comment"] != nil {
				assert.Equal(t, responseMap["No_comment"], "No Comment Found")
			}
		}
	}
}

func TestDeleteComment(t *testing.T) {

	gin.SetMode(gin.TestMode)

	var secondUserEmail, secondUserPassword string
	// var secondUserID uint32
	var secondCommentID uint64

	err := refreshUserPostAndCommentTable()
	if err != nil {
		log.Fatal(err)
	}
	_, users, comments, err := seedUsersPostsAndComments()
	if err != nil {
		log.Fatalf("Cannot seed tables %v\n", err)
	}
	// Get only the second user
	for _, user := range users {
		if user.ID == 1 {
			continue
		}
		// secondUserID = user.ID
		secondUserEmail = user.Email
		secondUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	// Get only the second comment
	for _, comment := range comments {
		if comment.ID == 1 {
			continue
		}
		secondCommentID = comment.ID
	}

	//Login the user and get the authentication token
	tokenInterface, err := server.SignIn(secondUserEmail, secondUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	token := tokenInterface["token"] //get only the token
	tokenString := fmt.Sprintf("Bearer %v", token)

	commentsSample := []struct {
		commentID      string
		usersLength    int
		tokenGiven     string
		commentsLength int
		statusCode     int
	}{
		{
			commentID:  strconv.Itoa(int(secondCommentID)),
			statusCode: 200,
			tokenGiven: tokenString,
		},
		{
			//an id that does not exist
			commentID:  strconv.Itoa(12322),
			statusCode: 404,
			tokenGiven: tokenString,
		},
		{
			//When the user is not authenticated
			commentID:  strconv.Itoa(int(secondCommentID)),
			statusCode: 401,
			tokenGiven: "",
		},
		{
			//When wrong token is passed
			commentID:  strconv.Itoa(int(secondCommentID)),
			statusCode: 401,
			tokenGiven: "this is a wrong token",
		},
		{
			// When id passed is invalid
			commentID:  "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range commentsSample {

		r := gin.Default()
		r.DELETE("/comments/:id", server.DeleteComment)
		req, err := http.NewRequest(http.MethodDelete, "/comments/"+v.commentID, nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		req.Header.Set("Authorization", v.tokenGiven)
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
		if err != nil {
			t.Errorf("Cannot convert to json here: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			responseMap := responseInterface["response"]
			assert.Equal(t, responseMap, "Comment deleted")
		}
		if v.statusCode == 400 || v.statusCode == 401 || v.statusCode == 404 {
			responseMap := responseInterface["error"].(map[string]interface{})
			if responseMap["Invalid_request"] != nil {
				assert.Equal(t, responseMap["Invalid_request"], "Invalid Request")
			}
			if responseMap["Unauthorized"] != nil {
				assert.Equal(t, responseMap["Unauthorized"], "Unauthorized")
			}
			if responseMap["No_comment"] != nil {
				assert.Equal(t, responseMap["No_comment"], "No Comment Found")
			}
		}
	}
}
