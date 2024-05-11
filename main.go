package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
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

	// userMock := getUserById(db, 4)
	// userMock.FirstName = "Artist"
	// userMock.LastName = "Extraordinaire"
	// userMock.Username = "artist_extraordinaire"
	// updateUser(db, userMock)
	app := fiber.New()

	app.Get("/users", func(c *fiber.Ctx) error {
		return c.JSON(getAllUsers(db))
	})
	app.Get("/user/:id", func(c *fiber.Ctx) error {
		userId, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		userData := getUserById(db, userId)
		return c.JSON(userData)
	})
	app.Post("/users/create", func(c *fiber.Ctx) error {
		newUser := new(User)
		if err := c.BodyParser(newUser); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		err = createUser(db, newUser)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		return c.JSON(fiber.Map{
			"message": "Create user successfully",
		})

	})
	app.Put("/user/:id", func(c *fiber.Ctx) error {
		userId, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		newUser := new(User)
		if err := c.BodyParser(newUser); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		newUser.ID = uint(userId)
		err = updateUser(db, newUser)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		return c.JSON(fiber.Map{
			"message": "Update user successfully",
		})

	})
	app.Listen(":8080")
}
