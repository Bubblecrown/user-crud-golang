// คำสั่งจัดการ database
// structure ในการสร้าง database

package main

import (
	"log"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Tier      uint   `json:"tier"`
}

func getAllUsers(db *gorm.DB) []User {
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		log.Fatalf("Error creating user: %v", result.Error)
	}
	return users
}

func getUserById(db *gorm.DB, id int) *User {
	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		log.Fatalf("Error creating user: %v", result.Error)
	}
	return &user
}

func createUser(db *gorm.DB, user *User) error {
	result := db.Create(user)
	if result.Error != nil {
		log.Fatalf("Error creating user: %v", result.Error)
	}
	return nil
}

func updateUser(db *gorm.DB, user *User) error {
	result := db.Model(&user).Updates(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func deleteUser(db *gorm.DB, id int) error {
	var user User
	// Soft delete protects against accidental data loss. In contrast, hard delete means that if data is removed, it's gone forever.

	// permanently delete
	// result := db.Unscoped().Delete(&user, id)

	// soft delete
	result := db.Delete(&user, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func searchUserByFirstName(db *gorm.DB, firstName string) []User {
	var user []User
	result := db.Where("first_name = ?", firstName).Order("tier").Find(&user)
	if result.Error != nil {
		log.Fatalf("Error seaching user: %v", result.Error)
	}
	return user
}
