package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestForgotPassword(t *testing.T) {

	gin.SetMode(gin.TestMode)

	err := refreshUserAndResetPasswordTable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedOneUser()
	if err != nil {
		log.Fatal(err)
	}
	samples := []struct {
		id         string
		inputJSON  string
		statusCode int
	}{
		{
			// When the user input invalid email:
			inputJSON:  `{"email": "petgmail.com"}`,
			statusCode: 422,
		},
		{
			// When the email given dont exist in our database:
			inputJSON:  `{"email": "raman@example.com"}`,
			statusCode: 422,
		},
		{
			// When the email field is empty:
			inputJSON:  `{"email": ""}`,
			statusCode: 422,
		},
		{
			// When the number or any other input that is not string:
			inputJSON:  `{"email": 123}`,
			statusCode: 422,
		},
		// {
		// We comment the passing test. Why?
		//It will actually send the mail
		// This is not ideal in a testing environment
		// You can mock the process using Interface, or if you have a better idea,
		// You can raise a PR.

		// inputJSON:  `{"email": "pet@gmail.com"}`, //the seeded user
		// statusCode: 200,
		// },
	}
	for _, v := range samples {
		r := gin.Default()
		r.POST("/password/forgot", server.ForgotPassword)
		req, err := http.NewRequest(http.MethodPost, "/password/forgot", bytes.NewBufferString(v.inputJSON))
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

		// This is commented because, it not a good idea of sending real email while testing(it will also consume time)
		// if v.statusCode == 200 {
		// 	responseMap := responseInterface["response"]
		// 	assert.Equal(t, responseMap, "Success, Please click on the link provided in your email")
		// }
		if v.statusCode == 422 {
			responseMap := responseInterface["error"].(map[string]interface{})

			if responseMap["Invalid_email"] != nil {
				assert.Equal(t, responseMap["Invalid_email"], "Invalid Email")
			}
			if responseMap["No_email"] != nil {
				assert.Equal(t, responseMap["No_email"], "Sorry, we do not recognize this email")
			}
			if responseMap["Required_email"] != nil {
				assert.Equal(t, responseMap["Required_email"], "Required Email")
			}
			if responseMap["Unmarshal_error"] != nil {
				assert.Equal(t, responseMap["Unmarshal_error"], "Cannot unmarshal body")
			}
		}
	}
}

// func TestResetPassword(t *testing.T) {

// 	gin.SetMode(gin.TestMode)

// 	err := refreshUserTable()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	user, err := seedOneUser()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// The unhashed user password
// 	password := "password"

// 	//Login the user and get the authentication token
// 	tokenInterface, err := server.SignIn(user.Email, password)
// 	if err != nil {
// 		log.Fatalf("cannot login: %v\n", err)
// 	}
// 	token := tokenInterface["token"] //get only the token
// 	tokenString := fmt.Sprintf("Bearer %v", token)

// 	samples := []struct {
// 		id          string
// 		updateJSON  string
// 		statusCode  int
// 		username    string
// 		updateEmail string
// 		tokenGiven  string
// 	}{
// 		{
// 			id:          strconv.Itoa(int(user.ID)),
// 			updateJSON:  `{"email": "grand@gmail.com", "current_password": "password", "new_password": "newpassword"}`,
// 			statusCode:  200,
// 			username:    user.Username, //the username does not change, even if a new name is provided, it will be ignored
// 			updateEmail: user.Email,
// 			tokenGiven:  tokenString,
// 		},
// 	}

// 	for _, v := range samples {

// 		r := gin.Default()

// 		r.PUT("/users/:id", server.UpdateUser)
// 		req, err := http.NewRequest(http.MethodPut, "/users/"+v.id, bytes.NewBufferString(v.updateJSON))
// 		req.Header.Set("Authorization", v.tokenGiven)
// 		if err != nil {
// 			t.Errorf("this is the error: %v\n", err)
// 		}
// 		rr := httptest.NewRecorder()
// 		r.ServeHTTP(rr, req)

// 		responseInterface := make(map[string]interface{})
// 		err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
// 		if err != nil {
// 			t.Errorf("Cannot convert to json: %v", err)
// 		}

// 		assert.Equal(t, rr.Code, v.statusCode)

// 		if v.statusCode == 200 {
// 			//casting the interface to map:
// 			responseMap := responseInterface["response"].(map[string]interface{})
// 			assert.Equal(t, responseMap["email"], v.updateEmail)
// 			// assert.Equal(t, responseMap["username"], v.username)
// 		}

// 		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 {
// 			responseMap := responseInterface["error"].(map[string]interface{})

// 			fmt.Println("this is the response error: ", responseMap)

// 			if responseMap["Password_mismatch"] != nil {
// 				assert.Equal(t, responseMap["Password_mismatch"], "The password not correct")
// 			}
// 		}
// 	}
// }
