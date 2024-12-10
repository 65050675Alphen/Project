package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User Table
type User struct {
	gorm.Model
	Email      string `gorm:"unique"`
	Password   string `json:"password"`
	CustomerID uint
	Customer   Customer `gorm:"foreignKey:CustomerID"`
	Pets       []Pet    `gorm:"many2many:user_id;"`
}

func createUser(db *gorm.DB, user *User) error {
	hashedPassword, err :=
		bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
func loginUser(db *gorm.DB, user *User) (string, string, uint, error) {
	selectUser := new(User)
	// Correct the query to search by email
	result := db.Where("email = ?", user.Email).First(selectUser)

	if result.Error != nil {
		return "", "", 0, result.Error
	}

	// Verify the password
	err := bcrypt.CompareHashAndPassword(
		[]byte(selectUser.Password),
		[]byte(user.Password),
	)
	if err != nil {
		return "", "", 0, err
	}

	// Fetch the associated customer
	var customer Customer
	db.First(&customer, selectUser.CustomerID)

	env_err := godotenv.Load()
	if env_err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecretKey := os.Getenv("JWT_SECRET")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["user_id"] = selectUser.ID
	claims["customer_id"] = selectUser.CustomerID // เพิ่ม customer_id เข้าไปใน JWT claims
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", "", 0, err
	}

	return t, customer.CustomerName, selectUser.CustomerID, nil // Return token, customer name, and customerID
}

// getUser retrieves all user
func getUsers(db *gorm.DB, c *fiber.Ctx) error {
	var user []User
	db.Find(&user)
	return c.JSON(user)
}

// getUser retrieves a User by id
func getUser(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var user User
	db.First(&user, id)
	return c.JSON(user)
}

// updateUser updates a User by id
func updateUser(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	user := new(User)
	db.First(&user, id)
	if err := c.BodyParser(user); err != nil {
		return err
	}
	db.Save(&user)
	return c.JSON(user)
}

// deleteUser deletes a User by id
func deleteUser(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	db.Delete(&User{}, id)
	return c.SendString("Book successfully deleted")
}
