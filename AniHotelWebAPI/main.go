package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func authRequired(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	env_err := godotenv.Load()
	if env_err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecretKey := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})

	if err != nil || !token.Valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.Next()
}
func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")         // or the Docker service name if running in another container
	portStr := os.Getenv("DB_PORT")      // default PostgreSQL port
	user := os.Getenv("DB_USER")         // as defined in docker-compose.yml
	password := os.Getenv("DB_PASSWORD") // as defined in docker-compose.yml
	dbname := os.Getenv("DB_NAME")       // as defined in docker-compose.yml

	fmt.Println(host, portStr, user, password, dbname)
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port: %v", err)
	}

	// Configure your PostgreSQL database details here
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// New logger for detailed SQL logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the schema
	db.AutoMigrate(
		&Customer{},
		&User{},
		&Pet{},
	)
	// Setup Fiber
	app := fiber.New()

	//initialAPI
	initialCustomerAPI(db, app)
	initialUserAPI(db, app)
	initialPetAPI(db, app)
	// Start server
	log.Fatal(app.Listen(":8000"))
}

func initialPetAPI(db *gorm.DB, app *fiber.App) {
	// CRUD routes for pets with customer ID
	app.Post("/customers/:customerID/pets", func(c *fiber.Ctx) error {
		return createPetWithCustomerID(db, c)
	})

	app.Get("/customers/:customerID/pets", func(c *fiber.Ctx) error {
		return getPets(db, c)
	})

	app.Get("/customers/:customerID/pets/:id", func(c *fiber.Ctx) error {
		return getPet(db, c)
	})

	app.Put("/customers/:customerID/pets/:id", func(c *fiber.Ctx) error {
		return updatePet(db, c)
	})

	app.Delete("/customers/:customerID/pets/:id", func(c *fiber.Ctx) error {
		return deletePet(db, c)
	})
}

func initialCustomerAPI(db *gorm.DB, app *fiber.App) {

	// // CRUD routes
	app.Post("/customers", func(c *fiber.Ctx) error {
		return createCustomer(db, c)
	})

	app.Get("/customers", func(c *fiber.Ctx) error {
		return getCustomers(db, c)
	})
	app.Get("/customers/:id", func(c *fiber.Ctx) error {
		return getCustomer(db, c)
	})
	app.Put("/customers/:id", func(c *fiber.Ctx) error {
		return updateCustomer(db, c)
	})
	app.Delete("/customers/:id", func(c *fiber.Ctx) error {
		return deleteCustomer(db, c)
	})
}

func initialUserAPI(db *gorm.DB, app *fiber.App) {
	// app.Post("/User", func(c *fiber.Ctx) error {
	// 	return createUser(db, c)
	// })

	// //User API
	app.Post("/register", func(c *fiber.Ctx) error {
		user := new(User)

		if err := c.BodyParser(user); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		err := createUser(db, user)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.JSON(fiber.Map{
			"message": "Register Successful",
		})
	})

	app.Get("/users", func(c *fiber.Ctx) error {
		return getUsers(db, c)
	})

	app.Get("/users/:id", func(c *fiber.Ctx) error {
		return getUser(db, c)
	})

	app.Put("/users/:id", func(c *fiber.Ctx) error {
		return updateUser(db, c)
	})

	app.Delete("/users/:id", func(c *fiber.Ctx) error {
		return deleteUser(db, c)
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		// เรียก loginUser เพื่อให้ได้ token, customerName, และ customerID
		token, customerName, customerID, err := loginUser(db, user)

		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// Set cookie
		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(time.Hour * 72),
			HTTPOnly: true,
		})

		// ส่ง token, customerName และ customerID กลับไปใน response
		return c.JSON(fiber.Map{
			"token":        token,
			"customerName": customerName, // ส่ง customerName ใน response
			"customerID":   customerID,   // ส่ง customerID ใน response
		})
	})

}
