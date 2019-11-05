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

func TestResetPassword(t *testing.T) {

	gin.SetMode(gin.TestMode)

	err := refreshUserAndResetPasswordTable()
	if err != nil {
		log.Fatal(err)
	}
	// This is important when we want to update the user password
	_, err = seedOneUser()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedResetPassword()
	if err != nil {
		log.Fatal(err)
	}

	samples := []struct {
		inputJSON  string
		statusCode int
	}{
		{
			// When no token is passed:
			inputJSON:  `{"token": ""}`,
			statusCode: 422,
		},
		{
			// When the token is tampered with:
			inputJSON:  `{"token": "23423498398rwnef9sd8fjsdf"}`,
			statusCode: 422,
		},
		{
			// When passwords "new_password" and "retype_password" provided are not up to 6 characters:
			inputJSON:  `{"token": "awesometoken", "new_password": "pass", "retype_password":"pass"}`,
			statusCode: 422,
		},
		{
			// When the "new_password" is empty:
			inputJSON:  `{"token": "awesometoken", "new_password": "", "retype_password":"password"}`,
			statusCode: 422,
		},
		{
			// When the "retype_password" is empty:
			inputJSON:  `{"token": "awesometoken", "new_password": "password", "retype_password":""}`,
			statusCode: 422,
		},
		{
			// When the two password fields dont match
			inputJSON:  `{"token": "awesometoken", "new_password": "password", "retype_password":"newpassword"}`,
			statusCode: 422,
		},
		{
			// When the token and the password fields are correct, and the password updated
			inputJSON:  `{"token": "awesometoken", "new_password": "password", "retype_password":"password"}`,
			statusCode: 200,
		},
	}
	for _, v := range samples {
		r := gin.Default()
		r.POST("/password/reset", server.ResetPassword)
		req, err := http.NewRequest(http.MethodPost, "/password/reset", bytes.NewBufferString(v.inputJSON))
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
			responseMap := responseInterface["response"]
			assert.Equal(t, responseMap, "Success")
		}
		if v.statusCode == 422 {
			responseMap := responseInterface["error"].(map[string]interface{})

			if responseMap["Invalid_token"] != nil {
				assert.Equal(t, responseMap["Invalid_token"], "Invalid link. Try requesting again")
			}
			if responseMap["No_email"] != nil {
				assert.Equal(t, responseMap["No_email"], "Sorry, we do not recognize this email")
			}
			if responseMap["Invalid_Passwords"] != nil {
				assert.Equal(t, responseMap["Invalid_Passwords"], "Password should be atleast 6 characters")
			}
			if responseMap["Empty_passwords"] != nil {
				assert.Equal(t, responseMap["Empty_passwords"], "Please ensure both field are entered")
			}
			if responseMap["Password_unequal"] != nil {
				assert.Equal(t, responseMap["Password_unequal"], "Passwords provided do not match")
			}
		}
	}
}
