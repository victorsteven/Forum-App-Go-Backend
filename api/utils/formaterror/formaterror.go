package formaterror

import (
	"errors"
	"strings"
)

var errorMessages []string
var err error

func FormatError(errString string) []string {

	if strings.Contains(errString, "nickname") {
		err = errors.New("Nickname Already Taken")
		errorMessages = append(errorMessages, err.Error())
	}
	if strings.Contains(errString, "email") {
		err = errors.New("Email Already Taken")
		errorMessages = append(errorMessages, err.Error())
	}
	if strings.Contains(errString, "title") {
		err = errors.New("Title Already Taken")
		errorMessages = append(errorMessages, err.Error())
	}
	if strings.Contains(errString, "hashedPassword") {
		err = errors.New("Incorrect Password")
		errorMessages = append(errorMessages, err.Error())
	}
	if len(errorMessages) > 0 {
		return errorMessages
	}
	if len(errorMessages) == 0 {
		errorMessages = append(errorMessages, "Incorrect Details")
		return errorMessages
	}
	return nil
}
