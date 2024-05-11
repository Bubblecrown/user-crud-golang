package main

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Credential struct {
	gorm.Model
	Email    string `gorm:"unique"` // unique email address
	Password string
}

func createCredential(db *gorm.DB, credential *Credential) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credential.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	credential.Password = string(hashedPassword)
	result := db.Create(credential)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
