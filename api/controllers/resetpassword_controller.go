package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/victorsteven/fullstack/api/mailer"
	"github.com/victorsteven/fullstack/api/models"
	"github.com/victorsteven/fullstack/api/security"
	"io/ioutil"
	"net/http"
)

func (server *Server) ForgotPassword(c *gin.Context) {

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

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", requestBody["email"]).Take(&user).Error
	if err != nil {
		errList["No_Email"] = "The email does not exist"
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
		//formattedError := formaterror.FormatError(err.Error())
		//errList = formattedError
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err,
		})
		return
	}
	//c.JSON(http.StatusCreated, gin.H{
	//	"status":   http.StatusCreated,
	//	"response": userCreated,
	//})

	//Send welcome mail to the user:
	err = mailer.SendResetPassword(resetRecord)
	if err != nil {
		fmt.Printf("this is the sending mail error: %s\n", err)
	}
}