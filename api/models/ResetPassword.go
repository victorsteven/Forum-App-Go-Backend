package models

import (
	"github.com/jinzhu/gorm"
	"html"
	//"github.com/victorsteven/fullstack/api/security"
	"strings"
)

type ResetPassword struct {
	gorm.Model
	Email     string    `gorm:"size:100;not null;" json:"email"`
	Token  string    `gorm:"size:255;not null;" json:"token"`
}

func (resetPassword *ResetPassword) Prepare() {
	resetPassword.ID = 0
	resetPassword.Token = html.EscapeString(strings.TrimSpace(resetPassword.Token))
	resetPassword.Email = html.EscapeString(strings.TrimSpace(resetPassword.Email))
}

func (resetPassword *ResetPassword) SaveDatails(db *gorm.DB) (*ResetPassword, error) {
	var err error
	err = db.Debug().Create(&resetPassword).Error
	if err != nil {
		return &ResetPassword{}, err
	}
	return resetPassword, nil
}