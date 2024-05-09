// คำสั่งจัดการ database
// structure ในการสร้าง database

package main

import (
	"fmt"
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

func getUserById(db *gorm.DB, id int) *User {
	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		log.Fatalf("Error creating user: %v", result.Error)
	}
	return &user
}

func createUser(db *gorm.DB, user *User) {
	result := db.Create(user)
	if result.Error != nil {
		log.Fatalf("Error creating user: %v", result.Error)
	}
	fmt.Printf("User created successfully")
}

func updateUser(db *gorm.DB, user *User) {
	result := db.Save(&user)
	if result.Error != nil {
		log.Fatalf("Error updating user: %v", result.Error)
	}
	fmt.Printf("Update user successfully")
}

func deleteUser(db *gorm.DB, id int) {
	var user User
	// don't soft deletes the user
	// result := db.Unscoped().Delete(&user, id)
	result := db.Delete(&user, id)
	if result.Error != nil {
		log.Fatalf("Error deleting user: %v", result.Error)
	}
	fmt.Printf("Delete user successfully")
}

func searchUserByFirstName(db *gorm.DB, firstName string) []User {
	var user []User
	result := db.Where("first_name = ?", firstName).Order("tier").Find(&user)
	if result.Error != nil {
		log.Fatalf("Error seaching user: %v", result.Error)
	}
	return user
}
