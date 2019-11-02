package tests

import (
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

func TestLikePost(t *testing.T) {

	err := refreshUserPostAndLikeTable()
	if err != nil {
		log.Fatal(err)
	}
	post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	// Note: the value of the user password before it was hashed is "password". so:
	password := "password"
	tokenInterface, err := server.SignIn(post.Author.Email, password) //get the auth user email from the post
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	token := tokenInterface["token"] //get only the token
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		postID     string
		statusCode int
		tokenGiven string
	}{
		{
			postID:     strconv.Itoa(int(post.ID)),
			statusCode: 201,
			tokenGiven: tokenString,
		},
		{
			postID:     strconv.Itoa(int(post.ID)),
			statusCode: 401,
			tokenGiven: "",
		},
		{
			postID:     strconv.Itoa(int(post.ID)),
			statusCode: 401,
			tokenGiven: "This is an incorrect token",
		},
	}

	for _, v := range samples {

		gin.SetMode(gin.TestMode)

		r := gin.Default()

		r.POST("/likes/:id", server.LikePost)
		req, err := http.NewRequest(http.MethodPost, "/likes/"+v.postID, nil)
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
			fmt.Println("this is the like response: ", responseMap)
			// assert.Equal(t, responseMap["title"], v.title)
			// assert.Equal(t, responseMap["content"], v.content)
		}

		// if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 {
		// 	responseMap := responseInterface["error"].(map[string]interface{})

		// 	if responseMap["Unauthorized"] != nil {
		// 		assert.Equal(t, responseMap["Unauthorized"], "Unauthorized")
		// 	}
		// 	if responseMap["Taken_title"] != nil {
		// 		assert.Equal(t, responseMap["Taken_title"], "Title Already Taken")
		// 	}
		// 	if responseMap["Required_title"] != nil {
		// 		assert.Equal(t, responseMap["Required_title"], "Required Title")
		// 	}
		// 	if responseMap["Required_content"] != nil {
		// 		assert.Equal(t, responseMap["Required_content"], "Required Content")
		// 	}
		// }
	}
}
