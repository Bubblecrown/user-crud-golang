package main

import (
	"time"
	"github.com/golang-jwt/jwt/v4"
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

func loginCredential(db *gorm.DB, credential *Credential) (string, error) {
	loginUser := new(Credential)
	result := db.Where("email=?", credential.Email).First(loginUser)
	if result.Error != nil {
		return "", result.Error
	}
	err := bcrypt.CompareHashAndPassword([]byte(loginUser.Password), []byte(credential.Password))
	if err != nil {
		return "", err
	}
	// Create JWT token
	jwtSecretKey := "secret"
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = loginUser.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}
	return t, nil
}
