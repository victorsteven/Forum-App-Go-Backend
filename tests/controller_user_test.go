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
	"gopkg.in/go-playground/assert.v1"
)

// func TestCreateUser(t *testing.T) {

// 	err := refreshUserTable()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	samples := []struct {
// 		inputJSON    string
// 		statusCode   int
// 		username     string
// 		email        string
// 		errorMessage string
// 	}{
// 		{
// 			inputJSON:    `{"username":"Pet", "email": "pet@gmail.com", "password": "password"}`,
// 			statusCode:   201,
// 			username:     "Pet",
// 			email:        "pet@gmail.com",
// 			errorMessage: "",
// 		},
// 		{
// 			inputJSON:    `{"username":"Frank", "email": "pet@gmail.com", "password": "password"}`,
// 			statusCode:   500,
// 			errorMessage: "Email Already Taken",
// 		},
// 		{
// 			inputJSON:    `{"username":"Pet", "email": "grand@gmail.com", "password": "password"}`,
// 			statusCode:   500,
// 			errorMessage: "Username Already Taken",
// 		},
// 		{
// 			inputJSON:    `{"username":"Kan", "email": "kangmail.com", "password": "password"}`,
// 			statusCode:   422,
// 			errorMessage: "Invalid Email",
// 		},
// 		{
// 			inputJSON:    `{"username": "", "email": "kan@gmail.com", "password": "password"}`,
// 			statusCode:   422,
// 			errorMessage: "Required Username",
// 		},
// 		{
// 			inputJSON:    `{"username": "Kan", "email": "", "password": "password"}`,
// 			statusCode:   422,
// 			errorMessage: "Required Email",
// 		},
// 		{
// 			inputJSON:    `{"username": "Kan", "email": "kan@gmail.com", "password": ""}`,
// 			statusCode:   422,
// 			errorMessage: "Required Password",
// 		},
// 	}

// 	for _, v := range samples {

// 		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(v.inputJSON))
// 		if err != nil {
// 			t.Errorf("this is the error: %v", err)
// 		}
// 		rr := httptest.NewRecorder()
// 		handler := http.HandlerFunc(server.CreateUser)
// 		handler.ServeHTTP(rr, req)

// 		responseMap := make(map[string]interface{})
// 		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
// 		if err != nil {
// 			fmt.Printf("Cannot convert to json: %v", err)
// 		}
// 		assert.Equal(t, rr.Code, v.statusCode)
// 		if v.statusCode == 201 {
// 			assert.Equal(t, responseMap["username"], v.username)
// 			assert.Equal(t, responseMap["email"], v.email)
// 		}
// 		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
// 			assert.Equal(t, responseMap["error"], v.errorMessage)
// 		}
// 	}
// }

func TestGetUsers(t *testing.T) {

	// Switch to test mode so you don't get such noisy output
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

	// Switch to test mode so you don't get such noisy output
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
		id           string
		statusCode   int
		username     string
		email        string
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(user.ID)),
			statusCode: 200,
			username:   user.Username,
			email:      user.Email,
		},
		{
			id:           "unknwon",
			statusCode:   400,
			errorMessage: "Invalid Request",
		},
	}
	for _, v := range userSample {
		req, _ := http.NewRequest("GET", "/users/"+v.id, nil)
		rr := httptest.NewRecorder()

		r := gin.Default()
		r.GET("/users/:id", server.GetUser)
		r.ServeHTTP(rr, req)

		userMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &userMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		theUser := userMap["response"]                  // Get the response from the payload
		userData, _ := theUser.(map[string]interface{}) //converting theUser to a map from a interface

		theError := userMap["error"]                      // Get the error from the payload
		errorData, _ := theError.(map[string]interface{}) //converting theError to a map from a interface

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, user.Username, userData["username"])
			assert.Equal(t, user.Email, userData["email"])
		}
		if v.statusCode == 400 {
			assert.Equal(t, v.errorMessage, errorData["invalid_request"])
		}
	}
}

func TestUpdateUser(t *testing.T) {

	// Switch to test mode so you don't get such noisy output
	gin.SetMode(gin.TestMode)

	var AuthEmail, AuthPassword string
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
		AuthPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		id             string
		updateJSON     string
		statusCode     int
		updateUsername string
		updateEmail    string
		tokenGiven     string
		errorMessage   string
	}{
		{
			// Convert int32 to int first before converting to string
			id:             strconv.Itoa(int(AuthID)),
			updateJSON:     `{"username":"Grand", "email": "grand@gmail.com", "password": "password"}`,
			statusCode:     200,
			updateUsername: "Grand",
			updateEmail:    "grand@gmail.com",
			tokenGiven:     tokenString,
			errorMessage:   "",
		},
		{
			// When password field is empty
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"username":"Woman", "email": "woman@gmail.com", "password": ""}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Password",
		},
		// {
		// 	// When no token was passed
		// 	id:           strconv.Itoa(int(AuthID)),
		// 	updateJSON:   `{"username":"Man", "email": "man@gmail.com", "password": "password"}`,
		// 	statusCode:   401,
		// 	tokenGiven:   "",
		// 	errorMessage: "Unauthorized",
		// },
		// {
		// 	// When incorrect token was passed
		// 	id:           strconv.Itoa(int(AuthID)),
		// 	updateJSON:   `{"username":"Woman", "email": "woman@gmail.com", "password": "password"}`,
		// 	statusCode:   401,
		// 	tokenGiven:   "This is incorrect token",
		// 	errorMessage: "Unauthorized",
		// },
		// {
		// 	// Remember "kenny@gmail.com" belongs to user 2
		// 	id:           strconv.Itoa(int(AuthID)),
		// 	updateJSON:   `{"username":"Frank", "email": "kenny@gmail.com", "password": "password"}`,
		// 	statusCode:   500,
		// 	tokenGiven:   tokenString,
		// 	errorMessage: "Email Already Taken",
		// },
		// {
		// 	// Remember "Kenny Morris" belongs to user 2
		// 	id:           strconv.Itoa(int(AuthID)),
		// 	updateJSON:   `{"username":"Kenny Morris", "email": "grand@gmail.com", "password": "password"}`,
		// 	statusCode:   500,
		// 	tokenGiven:   tokenString,
		// 	errorMessage: "Username Already Taken",
		// },
		// {
		// 	id:           strconv.Itoa(int(AuthID)),
		// 	updateJSON:   `{"username":"Kan", "email": "kangmail.com", "password": "password"}`,
		// 	statusCode:   422,
		// 	tokenGiven:   tokenString,
		// 	errorMessage: "Invalid Email",
		// },
		// {
		// 	id:           strconv.Itoa(int(AuthID)),
		// 	updateJSON:   `{"username": "", "email": "kan@gmail.com", "password": "password"}`,
		// 	statusCode:   422,
		// 	tokenGiven:   tokenString,
		// 	errorMessage: "Required Username",
		// },
		// {
		// 	id:           strconv.Itoa(int(AuthID)),
		// 	updateJSON:   `{"username": "Kan", "email": "", "password": "password"}`,
		// 	statusCode:   422,
		// 	tokenGiven:   tokenString,
		// 	errorMessage: "Required Email",
		// },
		// {
		// 	id:         "unknwon",
		// 	tokenGiven: tokenString,
		// 	statusCode: 400,
		// },
		// {
		// 	// When user 2 is using user 1 token
		// 	id:           strconv.Itoa(int(2)),
		// 	updateJSON:   `{"username": "Mike", "email": "mike@gmail.com", "password": "password"}`,
		// 	tokenGiven:   tokenString,
		// 	statusCode:   401,
		// 	errorMessage: "Unauthorized",
		// },
	}

	for _, v := range samples {
		req, err := http.NewRequest("GET", "/users/"+v.id, bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()

		r := gin.Default()
		r.GET("/users/:id", server.UpdateUser)
		req.Header.Set("Authorization", v.tokenGiven)
		r.ServeHTTP(rr, req)

		// fmt.Printf("This si the json: %v\n", rr.Body.String())

		payload := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &payload)
		if err != nil {
			t.Errorf("Cannot convert to json: %v\n", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		theUser := payload["response"] // Get the response from the payload
		if theUser != nil {
			userData, _ := theUser.(map[string]interface{}) //converting theUser to a mp from a interface
			assert.Equal(t, userData["username"], v.updateUsername)
			assert.Equal(t, userData["email"], v.updateEmail)
		}

		theError := payload["error"] // Get the error from the payload
		if theError != nil {
			errorData, _ := theError.(map[string]interface{}) //converting theUser to a mp from a interface
			if errorData["required_password"] != nil {
				assert.Equal(t, v.errorMessage, errorData["required_password"])
			}
			if errorData["unauthorized"] != nil {
				assert.Equal(t, v.errorMessage, errorData["unauthorized"])
			}
			// fmt.Printf("this is the error: %v", errorData)
			// errorString := fmt.Sprintf("%v", theError) //Convert the interface to string
			// replacer := strings.NewReplacer("[", "", "]", "") //remove square brackets
			// errorString = replacer.Replace(errorString)

			// fmt.Printf("this is the error: %v", theError)
			// assert.Equal(t, errorString, v.errorMessage)
		}
	}
}

// func TestDeleteUser(t *testing.T) {

// 	var AuthEmail, AuthPassword string
// 	var AuthID uint32

// 	err := refreshUserTable()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	users, err := seedUsers() //we need atleast two users to properly check the update
// 	if err != nil {
// 		log.Fatalf("Error seeding user: %v\n", err)
// 	}
// 	// Get only the first and log him in
// 	for _, user := range users {
// 		if user.ID == 2 {
// 			continue
// 		}
// 		AuthID = user.ID
// 		AuthEmail = user.Email
// 		AuthPassword = "password" ////Note the password in the database is already hashed, we want unhashed
// 	}

// 	//Login the user and get the authentication token
// 	token, err := server.SignIn(AuthEmail, AuthPassword)
// 	if err != nil {
// 		log.Fatalf("cannot login: %v\n", err)
// 	}
// 	tokenString := fmt.Sprintf("Bearer %v", token)

// 	userSample := []struct {
// 		id           string
// 		tokenGiven   string
// 		statusCode   int
// 		errorMessage string
// 	}{
// 		{
// 			// Convert int32 to int first before converting to string
// 			id:           strconv.Itoa(int(AuthID)),
// 			tokenGiven:   tokenString,
// 			statusCode:   204,
// 			errorMessage: "",
// 		},
// 		{
// 			// When no token is given
// 			id:           strconv.Itoa(int(AuthID)),
// 			tokenGiven:   "",
// 			statusCode:   401,
// 			errorMessage: "Unauthorized",
// 		},
// 		{
// 			// When incorrect token is given
// 			id:           strconv.Itoa(int(AuthID)),
// 			tokenGiven:   "This is an incorrect token",
// 			statusCode:   401,
// 			errorMessage: "Unauthorized",
// 		},
// 		{
// 			id:         "unknwon",
// 			tokenGiven: tokenString,
// 			statusCode: 400,
// 		},
// 		{
// 			// User 2 trying to use User 1 token
// 			id:           strconv.Itoa(int(2)),
// 			tokenGiven:   tokenString,
// 			statusCode:   401,
// 			errorMessage: "Unauthorized",
// 		},
// 	}

// 	for _, v := range userSample {

// 		req, err := http.NewRequest("GET", "/users", nil)
// 		if err != nil {
// 			t.Errorf("This is the error: %v\n", err)
// 		}
// 		req = mux.SetURLVars(req, map[string]string{"id": v.id})
// 		rr := httptest.NewRecorder()
// 		handler := http.HandlerFunc(server.DeleteUser)

// 		req.Header.Set("Authorization", v.tokenGiven)

// 		handler.ServeHTTP(rr, req)
// 		assert.Equal(t, rr.Code, v.statusCode)

// 		if v.statusCode == 401 && v.errorMessage != "" {
// 			responseMap := make(map[string]interface{})
// 			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
// 			if err != nil {
// 				t.Errorf("Cannot convert to json: %v", err)
// 			}
// 			assert.Equal(t, responseMap["error"], v.errorMessage)
// 		}
// 	}
// }
