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
	db.AutoMigrate(&User{}, &Credential{})

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
	app.Delete("/user/:id", func(c *fiber.Ctx) error {
		userId, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		err = deleteUser(db, userId)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		return c.JSON(fiber.Map{
			"message": "Delete user successfully",
		})

	})

	app.Post("/register", func(c *fiber.Ctx) error {
		newCreateUser := new(Credential)
		if err := c.BodyParser(newCreateUser); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		err = createCredential(db, newCreateUser)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		return c.JSON(fiber.Map{
			"message": "Register user successfully",
		})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		newLoginUser := new(Credential)
		if err := c.BodyParser(newLoginUser); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		token, err := loginCredential(db, newLoginUser)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		// Set cookie
		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(time.Hour * 72),
			HTTPOnly: true,
		})
		return c.JSON(fiber.Map{
			"message": "Login successfully",
		})
	})
	app.Listen(":8080")
}
