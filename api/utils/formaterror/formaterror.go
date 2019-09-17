package formaterror

import (
	"errors"
	"strings"
)

var errorMessages = make(map[string]string)

var err error

func FormatError(errString string) map[string]string {

	if strings.Contains(errString, "nickname") {
		err = errors.New("Nickname Already Taken")
		errorMessages["taken_nickname"] = err.Error()

	}
	if strings.Contains(errString, "email") {
		err = errors.New("Email Already Taken")
		errorMessages["taken_nickname"] = err.Error()

	}
	if strings.Contains(errString, "title") {
		err = errors.New("Title Already Taken")
		errorMessages["taken_title"] = err.Error()

	}
	if strings.Contains(errString, "hashedPassword") {
		err = errors.New("Incorrect Password")
		errorMessages["incorrect_password"] = err.Error()

	}
	if len(errorMessages) > 0 {
		return errorMessages
	}
	if len(errorMessages) == 0 {
		errorMessages["incorrect_details"] = "Incorrect Details"
		return errorMessages
	}
	return nil
}
