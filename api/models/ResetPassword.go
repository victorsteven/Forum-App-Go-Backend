package models

import (
	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	//"github.com/victorsteven/fullstack/api/security"
	"strings"
	"errors"
	"html"
)

type ResetPassword struct {
	gorm.Model
	Email     string    `gorm:"size:100;not null;" json:"email"`
	Token  string    `gorm:"size:255;not null;" json:"token"`
}


//func (u *User) BeforeSave() error {
//	hashedPassword, err := security.Hash(u.Password)
//	if err != nil {
//		return err
//	}
//	u.Password = string(hashedPassword)
//	return nil
//}


func (resetPassword *ResetPassword) Prepare() {
	resetPassword.ID = 0
	resetPassword.Token = html.EscapeString(strings.TrimSpace(resetPassword.Token))
	resetPassword.Email = html.EscapeString(strings.TrimSpace(resetPassword.Email))
}

func (resetPassword *ResetPassword) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	switch strings.ToLower(action) {

	case "forgotpassword":
		if resetPassword.Email == "" {
			err = errors.New("Required Email")
			errorMessages["Required_email"] = err.Error()
		}
		if resetPassword.Email != "" {
			if err = checkmail.ValidateFormat(resetPassword.Email); err != nil {
				err = errors.New("Invalid Email")
				errorMessages["Invalid_email"] = err.Error()
			}
		}

	default:
		if resetPassword.Email == "" {
			err = errors.New("Required Email")
			errorMessages["Required_email"] = err.Error()
		}
		if resetPassword.Email != "" {
			if err = checkmail.ValidateFormat(resetPassword.Email); err != nil {
				err = errors.New("Invalid Email")
				errorMessages["Invalid_email"] = err.Error()
			}
		}
	}
	return errorMessages
}

func (resetPassword *ResetPassword) SaveDatails(db *gorm.DB) (*ResetPassword, error) {

	var err error
	err = db.Debug().Create(&resetPassword).Error
	if err != nil {
		return &ResetPassword{}, err
	}
	return resetPassword, nil
}