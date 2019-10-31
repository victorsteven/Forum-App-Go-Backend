package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/badoux/checkmail"
	"github.com/gin-gonic/gin"
	"github.com/victorsteven/fullstack/api/mailer"
	"github.com/victorsteven/fullstack/api/models"
	"github.com/victorsteven/fullstack/api/security"
	"github.com/victorsteven/fullstack/api/utils/formaterror"
	"io/ioutil"
	"net/http"
)

func (server *Server) ForgotPassword(c *gin.Context) {
	//remove any possible error, because the frontend dont reload
	errList = map[string]string{}

	// Start processing the request
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	requestBody := map[string]string{}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	user := models.User{}
	resetPassword := models.ResetPassword{}

	resetPassword.Prepare()

	fmt.Printf("this is the user email: %s\n", requestBody["email"])

	if requestBody["email"] == "" {
		errList["Required_email"] = "Required Email"
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"status": http.StatusUnprocessableEntity,
				"error":  errList,
			})
		return
	}
	if requestBody["email"] != "" {
		if err = checkmail.ValidateFormat(requestBody["email"]); err != nil {
			errList["Invalid_email"] = "Invalid Email"
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"status": http.StatusUnprocessableEntity,
				"error":  errList,
			})
			return
		}
	}
	err = server.DB.Debug().Model(models.User{}).Where("email = ?", requestBody["email"]).Take(&user).Error
	if err != nil {
		errList["No_email"] = "Sorry, we do not recognize this email"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	//generate the token:
	token := security.TokenHash(requestBody["email"])
	resetPassword.Email = requestBody["email"]
	resetPassword.Token = token

	resetRecord, err := resetPassword.SaveDatails(server.DB)
	if err != nil {
		errList = formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	//Send welcome mail to the user:
	err = mailer.SendResetPassword(resetRecord)
	if err != nil {
		fmt.Printf("this is the sending mail error: %s\n", err)
	}
}


func (server *Server) ResetPassword(c *gin.Context) {
	//remove any possible error, because the frontend dont reload
	errList = map[string]string{}

	// Start processing the request
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	requestBody := map[string]string{}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	user := models.User{}
	resetPassword := models.ResetPassword{}

	err = server.DB.Debug().Model(models.ResetPassword{}).Where("token = ?", requestBody["token"]).Take(&resetPassword).Error
	if err != nil {
		errList["Invalid_token"] = "Invalid link. Try requesting again"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", resetPassword.Email).Take(&user).Error
	if err != nil {
		errList["No_email"] = "Sorry, we do not recognize this email"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	if requestBody["new_password"] == "" || requestBody["retype_password"] == ""  {
		errList["Empty_passwords"] = "Please ensure both field are entered"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	if requestBody["new_password"] != "" && requestBody["retype_password"] != ""  {
		//Also check if the new password
		if len(requestBody["new_password"]) < 6 || len(requestBody["retype_password"]) < 6  {
			errList["Invalid_Passwords"] = "Password should be atleast 6 characters"
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"status": http.StatusUnprocessableEntity,
				"error":  errList,
			})
			return
		}
		if requestBody["new_password"] != requestBody["retype_password"] {
			errList["Password_unequal"] = "Passwords provided do not match"
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"status": http.StatusUnprocessableEntity,
				"error":  errList,
			})
			return
		}
		//Note this password will be hashed before it is saved in the model
		user.Password = requestBody["new_password"]
		user.Email = resetPassword.Email

		//update the password
		user.Prepare()
		err := user.UpdatePassword(server.DB)
		if err != nil {
			errList["Cannot_save"] = "Cannot Save, Pls try again later"
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  errList,
			})
			return
		}
		//Delete the token record so is not used again:
		_, err = resetPassword.DeleteDatails(server.DB)
		if err != nil {
			errList["Cannot_delete"] = "Cannot Delete record, Pls try again later"
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  errList,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   http.StatusOK,
			"response": "Success",
		})
	}
}