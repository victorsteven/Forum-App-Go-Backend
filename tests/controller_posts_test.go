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

func TestCreatePost(t *testing.T) {

	gin.SetMode(gin.TestMode)

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	// Note: the value of the user password before it was hashed is "password". so:
	password := "password"
	tokenInterface, err := server.SignIn(user.Email, password)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	token := tokenInterface["token"] //get only the token
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Note that the author id is obtained from the token, so we dont pass it
	samples := []struct {
		inputJSON  string
		statusCode int
		title      string
		content    string
		tokenGiven string
	}{
		{
			inputJSON:  `{"title":"The title", "content": "the content"}`,
			statusCode: 201,
			tokenGiven: tokenString,
			title:      "The title",
			content:    "the content",
		},
		{
			// When the post title already exist
			inputJSON:  `{"title":"The title", "content": "the content"}`,
			statusCode: 500,
			tokenGiven: tokenString,
		},
		{
			// When no token is passed
			inputJSON:  `{"title":"When no token is passed", "content": "the content"}`,
			statusCode: 401,
			tokenGiven: "",
		},
		{
			// When incorrect token is passed
			inputJSON:  `{"title":"When incorrect token is passed", "content": "the content"}`,
			statusCode: 401,
			tokenGiven: "This is an incorrect token",
		},
		{
			inputJSON:  `{"title": "", "content": "The content"}`,
			statusCode: 422,
			tokenGiven: tokenString,
		},
		{
			inputJSON:  `{"title": "This is a title", "content": ""}`,
			statusCode: 422,
			tokenGiven: tokenString,
		},
	}

	for _, v := range samples {

		r := gin.Default()

		r.POST("/posts", server.CreatePost)
		req, err := http.NewRequest(http.MethodPost, "/posts", bytes.NewBufferString(v.inputJSON))
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
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["content"], v.content)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 {
			responseMap := responseInterface["error"].(map[string]interface{})

			if responseMap["Unauthorized"] != nil {
				assert.Equal(t, responseMap["Unauthorized"], "Unauthorized")
			}
			if responseMap["Taken_title"] != nil {
				assert.Equal(t, responseMap["Taken_title"], "Title Already Taken")
			}
			if responseMap["Required_title"] != nil {
				assert.Equal(t, responseMap["Required_title"], "Required Title")
			}
			if responseMap["Required_content"] != nil {
				assert.Equal(t, responseMap["Required_content"], "Required Content")
			}
		}
	}
}

func TestGetPosts(t *testing.T) {

	gin.SetMode(gin.TestMode)

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = seedUsersAndPosts()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/posts", server.GetUsers)

	req, err := http.NewRequest(http.MethodGet, "/posts", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	postsInterface := make(map[string]interface{})

	err = json.Unmarshal([]byte(rr.Body.String()), &postsInterface)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}
	// This is so that we can get the length of the posts:
	thePosts := postsInterface["response"].([]interface{})
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(thePosts), 2)
}

func TestGetPostByID(t *testing.T) {

	gin.SetMode(gin.TestMode)

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}
	_, post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatal(err)
	}

	postSample := []struct {
		id         string
		statusCode int
		title      string
		content    string
		author_id  uint32
	}{
		{
			id:         strconv.Itoa(int(post.ID)),
			statusCode: 200,
			title:      post.Title,
			content:    post.Content,
			author_id:  post.AuthorID,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
		{
			id:         strconv.Itoa(12322), //an id that does not exist
			statusCode: 404,
		},
	}
	for _, v := range postSample {
		req, _ := http.NewRequest("GET", "/posts/"+v.id, nil)
		rr := httptest.NewRecorder()

		r := gin.Default()
		r.GET("/posts/:id", server.GetPost)
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			responseMap := responseInterface["response"].(map[string]interface{})
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["content"], v.content)
			assert.Equal(t, responseMap["author_id"], float64(v.author_id))
		}
		if v.statusCode == 400 || v.statusCode == 404 {
			responseMap := responseInterface["error"].(map[string]interface{})

			if responseMap["Invalid_request"] != nil {
				assert.Equal(t, responseMap["Invalid_request"], "Invalid Request")
			}
			if responseMap["No_post"] != nil {
				assert.Equal(t, responseMap["No_post"], "No Post Found")
			}
		}
	}
}

func TestUpdatePost(t *testing.T) {

	gin.SetMode(gin.TestMode)

	var PostUserEmail, PostUserPassword string
	// var AuthID uint32
	var AuthPostID uint64

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}
	users, posts, err := seedUsersAndPosts()
	if err != nil {
		log.Fatal(err)
	}
	// Get only the first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		PostUserEmail = user.Email
		PostUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	// Get only the first post
	for _, post := range posts {
		if post.ID == 2 {
			continue
		}
		AuthPostID = post.ID
	}
	//Login the user and get the authentication token
	tokenInterface, err := server.SignIn(PostUserEmail, PostUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	token := tokenInterface["token"] //get only the token
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		id         string
		updateJSON string
		statusCode int
		title      string
		content    string
		tokenGiven string
	}{
		{
			// Convert int64 to int first before converting to string
			id:         strconv.Itoa(int(AuthPostID)),
			updateJSON: `{"title":"The updated post", "content": "This is the updated content"}`,
			statusCode: 200,
			title:      "The updated post",
			content:    "This is the updated content",
			tokenGiven: tokenString,
		},
		{
			// When no token is provided
			id:         strconv.Itoa(int(AuthPostID)),
			updateJSON: `{"title":"This is still another title", "content": "This is the updated content"}`,
			tokenGiven: "",
			statusCode: 401,
		},
		{
			// When incorrect token is provided
			id:         strconv.Itoa(int(AuthPostID)),
			updateJSON: `{"title":"This is still another title", "content": "This is the updated content"}`,
			tokenGiven: "this is an incorrect token",
			statusCode: 401,
		},
		{
			//Note: "Title 2" belongs to post 2, and title must be unique
			id:         strconv.Itoa(int(AuthPostID)),
			updateJSON: `{"title":"Title 2", "content": "This is the updated content"}`,
			statusCode: 500,
			tokenGiven: tokenString,
		},
		{
			// When title is not given
			id:         strconv.Itoa(int(AuthPostID)),
			updateJSON: `{"title":"", "content": "This is the updated content"}`,
			statusCode: 422,
			tokenGiven: tokenString,
		},
		{
			// When content is not given
			id:         strconv.Itoa(int(AuthPostID)),
			updateJSON: `{"title":"Awesome title", "content": ""}`,
			statusCode: 422,
			tokenGiven: tokenString,
		},
		{
			// When invalid post id is given
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range samples {

		r := gin.Default()

		r.PUT("/posts/:id", server.UpdatePost)
		req, err := http.NewRequest(http.MethodPut, "/posts/"+v.id, bytes.NewBufferString(v.updateJSON))
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

		if v.statusCode == 200 {
			//casting the interface to map:
			responseMap := responseInterface["response"].(map[string]interface{})
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["content"], v.content)
		}
		if v.statusCode == 400 || v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 {
			responseMap := responseInterface["error"].(map[string]interface{})
			if responseMap["Unauthorized"] != nil {
				assert.Equal(t, responseMap["Unauthorized"], "Unauthorized")
			}
			if responseMap["Invalid_request"] != nil {
				assert.Equal(t, responseMap["Invalid_request"], "Invalid Request")
			}
			if responseMap["Taken_title"] != nil {
				assert.Equal(t, responseMap["Taken_title"], "Title Already Taken")
			}
			if responseMap["Required_title"] != nil {
				assert.Equal(t, responseMap["Required_title"], "Required Title")
			}
			if responseMap["Required_content"] != nil {
				assert.Equal(t, responseMap["Required_content"], "Required Content")
			}
		}
	}
}

func TestDeletePost(t *testing.T) {

	gin.SetMode(gin.TestMode)

	var PostUserEmail, PostUserPassword string
	// var AuthID uint32
	var AuthPostID uint64

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}
	users, posts, err := seedUsersAndPosts()
	if err != nil {
		log.Fatal(err)
	}
	// Get only the second user
	for _, user := range users {
		if user.ID == 1 {
			continue
		}
		PostUserEmail = user.Email
		PostUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	// Get only the second post
	for _, post := range posts {
		if post.ID == 1 {
			continue
		}
		AuthPostID = post.ID
	}
	//Login the user and get the authentication token
	tokenInterface, err := server.SignIn(PostUserEmail, PostUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	token := tokenInterface["token"] //get only the token
	tokenString := fmt.Sprintf("Bearer %v", token)

	postSample := []struct {
		id           string
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:         strconv.Itoa(int(AuthPostID)),
			tokenGiven: tokenString,
			statusCode: 200,
		},
		{
			// When empty token is passed
			id:         strconv.Itoa(int(AuthPostID)),
			tokenGiven: "",
			statusCode: 401,
		},
		{
			// When incorrect token is passed
			id:         strconv.Itoa(int(AuthPostID)),
			tokenGiven: "This is an incorrect token",
			statusCode: 401,
		},
		{
			id:         "unknwon",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(1)),
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range postSample {
		r := gin.Default()
		r.DELETE("/posts/:id", server.DeletePost)
		req, _ := http.NewRequest(http.MethodDelete, "/posts/"+v.id, nil)
		req.Header.Set("Authorization", v.tokenGiven)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})

		err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
		if err != nil {
			t.Errorf("Cannot convert to json here: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, responseInterface["response"], "Post deleted")
		}

		if v.statusCode == 400 || v.statusCode == 401 {
			responseMap := responseInterface["error"].(map[string]interface{})

			if responseMap["Invalid_request"] != nil {
				assert.Equal(t, responseMap["Invalid_request"], "Invalid Request")
			}
			if responseMap["Unauthorized"] != nil {
				assert.Equal(t, responseMap["Unauthorized"], "Unauthorized")
			}
		}
	}
}
