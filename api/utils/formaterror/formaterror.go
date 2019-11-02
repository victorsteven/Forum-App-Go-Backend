package formaterror

import (
	"errors"
	"strings"
)

var errorMessages = make(map[string]string)

var err error

func FormatError(errString string) map[string]string {

	if strings.Contains(errString, "username") {
		err = errors.New("Username Already Taken")
		errorMessages["Taken_username"] = err.Error()
	}

	if strings.Contains(errString, "email") {
		err = errors.New("Email Already Taken")
		errorMessages["Taken_email"] = err.Error()

	}
	if strings.Contains(errString, "title") {
		err = errors.New("Title Already Taken")
		errorMessages["Taken_title"] = err.Error()

	}
	if strings.Contains(errString, "hashedPassword") {
		err = errors.New("Incorrect Password")
		errorMessages["Incorrect_password"] = err.Error()
	}
	if strings.Contains(errString, "record not found") {
		err = errors.New("No Record Found")
		errorMessages["No_record"] = err.Error()
	}

	if len(errorMessages) > 0 {
		return errorMessages
	}

	if len(errorMessages) == 0 {
		errorMessages["Incorrect_details"] = "Incorrect Details"
		return errorMessages
	}

	return nil
}
