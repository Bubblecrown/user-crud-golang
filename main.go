package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "0.0.0.0"    // or the Docker service name if running in another container
	port     = 5432         // default PostgreSQL port
	user     = "myuser"     // as defined in docker-compose.yml
	password = "mypassword" // as defined in docker-compose.yml
	dbname   = "mydatabase" // as defined in docker-compose.yml
)

func main() {
	// connect database
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		// break or kill process
		panic("failed to connect to database")
	}
	// compare table in database and struct in project
	db.AutoMigrate(&User{})
	fmt.Println("Database migration completed!")

	// Create User
	// newUser := &User{Username: "jane.smith", FirstName: "Jane", LastName: "Smith", Tier: 2}
	// createUser(db, newUser)

	// Get All Users
	// users := getUserById(db, 2)
	// formattedJson, err := json.MarshalIndent(users, "", "  ")
	// if err != nil {
	// 	fmt.Println("Error formatting JSON:", err)
	// 	return
	// }
	// fmt.Println(string(formattedJson))

	// Update User
	// users.Username = "nature_lover"
	// updateUser(db, users)

	// Delete User
	// deleteUser(db, 1)

	// Search user by username
	userByUsername := searchUserByFirstName(db, "Jane")
	formattedJson, err := json.MarshalIndent(userByUsername, "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}
	fmt.Println(string(formattedJson))

}
