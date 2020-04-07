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

func TestCreateUser(t *testing.T) {

	gin.SetMode(gin.TestMode)

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	samples := []struct {
		inputJSON  string
		statusCode int
		username   string
		email      string
	}{
		{
			inputJSON:  `{"username":"Pet", "email": "pet@example.com", "password": "password"}`,
			statusCode: 201,
			username:   "Pet",
			email:      "pet@example.com",
		},
		{
			inputJSON:  `{"username":"Frank", "email": "pet@example.com", "password": "password"}`,
			statusCode: 500,
		},
		{
			inputJSON:  `{"username":"Pet", "email": "grand@example.com", "password": "password"}`,
			statusCode: 500,
		},
		{
			inputJSON:  `{"username":"Kan", "email": "kanexample.com", "password": "password"}`,
			statusCode: 422,
		},
		{
			inputJSON:  `{"username": "", "email": "kan@example.com", "password": "password"}`,
			statusCode: 422,
		},
		{
			inputJSON:  `{"username": "Kan", "email": "", "password": "password"}`,
			statusCode: 422,
		},
		{
			inputJSON:  `{"username": "Kan", "email": "kan@example.com", "password": ""}`,
			statusCode: 422,
		},
	}

	for _, v := range samples {

		r := gin.Default()
		r.POST("/users", server.CreateUser)
		req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(v.inputJSON))
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
			//casting the interface to map:
			responseMap := responseInterface["response"].(map[string]interface{})
			assert.Equal(t, responseMap["username"], v.username)
			assert.Equal(t, responseMap["email"], v.email)
		}
		if v.statusCode == 422 || v.statusCode == 500 {
			responseMap := responseInterface["error"].(map[string]interface{})

			if responseMap["Taken_email"] != nil {
				assert.Equal(t, responseMap["Taken_email"], "Email Already Taken")
			}
			if responseMap["Taken_username"] != nil {
				assert.Equal(t, responseMap["Taken_username"], "Username Already Taken")
			}
			if responseMap["Invalid_email"] != nil {
				assert.Equal(t, responseMap["Invalid_email"], "Invalid Email")
			}
			if responseMap["Required_username"] != nil {
				assert.Equal(t, responseMap["Required_username"], "Required Username")
			}
			if responseMap["Required_email"] != nil {
				assert.Equal(t, responseMap["Required_email"], "Required Email")
			}
			if responseMap["Required_password"] != nil {
				assert.Equal(t, responseMap["Required_password"], "Required Password")
			}
		}
	}
}

func TestGetUsers(t *testing.T) {

	gin.SetMode(gin.TestMode)

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/users", server.GetUsers)

	req, err := http.NewRequest(http.MethodGet, "/users", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	usersMap := make(map[string]interface{})

	err = json.Unmarshal([]byte(rr.Body.String()), &usersMap)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}
	// This is so that we can get the length of the users:
	theUsers := usersMap["response"].([]interface{})
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(theUsers), 2)
}

func TestGetUserByID(t *testing.T) {

	gin.SetMode(gin.TestMode)

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}
	userSample := []struct {
		id         string
		statusCode int
		username   string
		email      string
	}{
		{
			id:         strconv.Itoa(int(user.ID)),
			statusCode: 200,
			username:   user.Username,
			email:      user.Email,
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
	for _, v := range userSample {
		req, _ := http.NewRequest("GET", "/users/"+v.id, nil)
		rr := httptest.NewRecorder()

		r := gin.Default()
		r.GET("/users/:id", server.GetUser)
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			responseMap := responseInterface["response"].(map[string]interface{})
			assert.Equal(t, responseMap["username"], v.username)
			assert.Equal(t, responseMap["email"], v.email)
		}

		if v.statusCode == 400 || v.statusCode == 404 {
			responseMap := responseInterface["error"].(map[string]interface{})

			if responseMap["Invalid_request"] != nil {
				assert.Equal(t, responseMap["Invalid_request"], "Invalid Request")
			}
			if responseMap["No_user"] != nil {
				assert.Equal(t, responseMap["No_user"], "No User Found")
			}
		}
	}
}

func TestUpdateUser(t *testing.T) {

	gin.SetMode(gin.TestMode)

	var AuthEmail, AuthPassword, AuthUsername string
	var AuthID uint32

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	users, err := seedUsers() //we need atleast two users to properly check the update
	if err != nil {
		log.Fatalf("Error seeding user: %v\n", err)
	}
	// Get only the first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		AuthID = user.ID
		AuthEmail = user.Email
		AuthUsername = user.Username
		AuthPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	tokenInterface, err := server.SignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	token := tokenInterface["token"] //get only the token
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		id          string
		updateJSON  string
		statusCode  int
		username    string
		updateEmail string
		tokenGiven  string
	}{
		{
			// In this particular test case, we changed the user's password to "newpassword". Very important to note
			// Convert int32 to int first before converting to string
			id:          strconv.Itoa(int(AuthID)),
			updateJSON:  `{"email": "grand@example.com", "current_password": "password", "new_password": "newpassword"}`,
			statusCode:  200,
			username:    AuthUsername, //the username does not change, even if a new name is provided, it will be ignored
			updateEmail: "grand@example.com",
			tokenGiven:  tokenString,
		},
		{
			// An attempt to change the username, will not work, the old name is still retained.
			// Remember, the "current_password" is now "newpassword", changed in test 1
			id:          strconv.Itoa(int(AuthID)),
			updateJSON:  `{"username": "new_name", "email": "grand@example.com", "current_password": "newpassword", "new_password": "newpassword"}`,
			statusCode:  200,
			username:    AuthUsername, //irrespective of the username inputed above, the old one is still used
			updateEmail: "grand@example.com",
			tokenGiven:  tokenString,
		},
		{
			// The user can update only his email address
			id:          strconv.Itoa(int(AuthID)),
			updateJSON:  `{"email": "fred@example.com"}`,
			statusCode:  200,
			username:    AuthUsername,
			updateEmail: "fred@example.com",
			tokenGiven:  tokenString,
		},
		{
			id:          strconv.Itoa(int(AuthID)),
			updateJSON:  `{"email": "alex@example.com", "current_password": "", "new_password": ""}`,
			statusCode:  200,
			username:    AuthUsername,
			updateEmail: "alex@example.com",
			tokenGiven:  tokenString,
		},
		{
			// When password the "current_password" is given and does not match with the one in the database
			id:          strconv.Itoa(int(AuthID)),
			updateJSON:  `{"email": "alex@example.com", "current_password": "wrongpassword", "new_password": "password"}`,
			statusCode:  422,
			updateEmail: "alex@example.com",
			tokenGiven:  tokenString,
		},
		{
			// When password the "current_password" is correct but the "new_password" field is not given
			id:          strconv.Itoa(int(AuthID)),
			updateJSON:  `{"email": "alex@example.com", "current_password": "newpassword", "new_password": ""}`,
			statusCode:  422,
			updateEmail: "alex@example.com",
			tokenGiven:  tokenString,
		},
		{
			// When password the "current_password" is correct but the "new_password" field is not up to 6 characters
			id:          strconv.Itoa(int(AuthID)),
			updateJSON:  `{"email": "alex@example.com", "current_password": "newpassword", "new_password": "pass"}`,
			statusCode:  422,
			updateEmail: "alex@example.com",
			tokenGiven:  tokenString,
		},
		{
			// When no token was passed (when the user is not authenticated)
			id:         strconv.Itoa(int(AuthID)),
			updateJSON: `{"email": "man@example.com", "current_password": "newpassword", "new_password": "password"}`,
			statusCode: 401,
			tokenGiven: "",
		},
		{
			// When incorrect token was passed
			id:         strconv.Itoa(int(AuthID)),
			updateJSON: `{"email": "man@example.com", "current_password": "newpassword", "new_password": "password"}`,
			statusCode: 401,
			tokenGiven: "This is incorrect token",
		},
		{
			// Remember "kenny@example.com" belongs to user 2, so, user 1 cannot use some else email that is in our database
			id:         strconv.Itoa(int(AuthID)),
			updateJSON: `{"email": "kenny@example.com", "current_password": "newpassword", "new_password": "password"}`,
			statusCode: 500,
			tokenGiven: tokenString,
		},
		{
			// When the email provided is invalid
			id:         strconv.Itoa(int(AuthID)),
			updateJSON: `{"email": "notexample.com", "current_password": "newpassword", "new_password": "password"}`,
			statusCode: 422,
			tokenGiven: tokenString,
		},
		{
			// If the email field is empty
			id:         strconv.Itoa(int(AuthID)),
			updateJSON: `{"email": "", "current_password": "newpassword", "new_password": "password"}`,
			statusCode: 422,
			tokenGiven: tokenString,
		},
		{
			// when invalid is provided
			id:         "unknwon",
			tokenGiven: tokenString,
			statusCode: 400,
		},
	}

	for _, v := range samples {

		r := gin.Default()

		r.PUT("/users/:id", server.UpdateUser)
		req, err := http.NewRequest(http.MethodPut, "/users/"+v.id, bytes.NewBufferString(v.updateJSON))
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
			assert.Equal(t, responseMap["email"], v.updateEmail)
			// assert.Equal(t, responseMap["username"], v.username)
		}

		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 {
			responseMap := responseInterface["error"].(map[string]interface{})

			fmt.Println("this is the response error: ", responseMap)

			if responseMap["Password_mismatch"] != nil {
				assert.Equal(t, responseMap["Password_mismatch"], "The password not correct")
			}
			if responseMap["Empty_new"] != nil {
				assert.Equal(t, responseMap["Empty_new"], "Please Provide new password")
			}
			if responseMap["Empty_current"] != nil {
				assert.Equal(t, responseMap["Empty_current"], "Please Provide current password")
			}
			if responseMap["Invalid_password"] != nil {
				assert.Equal(t, responseMap["Invalid_password"], "Password should be atleast 6 characters")
			}
			if responseMap["Unauthorized"] != nil {
				assert.Equal(t, responseMap["Unauthorized"], "Unauthorized")
			}
			if responseMap["Taken_email"] != nil {
				assert.Equal(t, responseMap["Taken_email"], "Email Already Taken")
			}
			if responseMap["Invalid_email"] != nil {
				assert.Equal(t, responseMap["Invalid_email"], "Invalid Email")
			}
			if responseMap["Required_email"] != nil {
				assert.Equal(t, responseMap["Required_email"], "Required Email")
			}
			if responseMap["Invalid_request"] != nil {
				assert.Equal(t, responseMap["Invalid_request"], "Invalid Request")
			}
		}
	}
}

func TestDeleteUser(t *testing.T) {

	gin.SetMode(gin.TestMode)

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}
	// Note: the value of the user password before it was hashed is "password". so:
	password := "password"
	tokenInterface, err := server.SignIn(user.Email, password)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	token := tokenInterface["token"] //get only the token
	tokenString := fmt.Sprintf("Bearer %v", token)

	userSample := []struct {
		id         string
		tokenGiven string
		statusCode int
	}{
		{
			// Convert int32 to int first before converting to string
			id:         strconv.Itoa(int(user.ID)),
			tokenGiven: tokenString,
			statusCode: 200,
		},
		{
			// When no token is given
			id:         strconv.Itoa(int(user.ID)),
			tokenGiven: "",
			statusCode: 401,
		},
		{
			// When incorrect token is given
			id:         strconv.Itoa(int(user.ID)),
			tokenGiven: "This is an incorrect token",
			statusCode: 401,
		},
		{
			// When bad request data is given:
			id:         "unknwon",
			tokenGiven: tokenString,
			statusCode: 400,
		},
	}
	for _, v := range userSample {

		r := gin.Default()
		r.DELETE("/users/:id", server.DeleteUser)
		req, _ := http.NewRequest(http.MethodDelete, "/users/"+v.id, nil)
		req.Header.Set("Authorization", v.tokenGiven)

		fmt.Println("FORM REQUEST: ", req.Header)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, responseInterface["response"], "User deleted")
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
