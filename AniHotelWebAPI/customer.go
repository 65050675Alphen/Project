package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Customer
type Customer struct {
	gorm.Model
	CustomerName  string `json:"name"`
	Description   string `json:"description"`
	CustomerPhone string `json:"phone"`
	Pets          []Pet  `json:"pets" gorm:"foreignKey:CustomerID"`
}

// createCustomer create a customer
func createCustomer(db *gorm.DB, c *fiber.Ctx) error {
	customer := new(Customer)
	if err := c.BodyParser(customer); err != nil {
		return err
	}
	db.Create(&customer)
	return c.JSON(customer)
}

// getCustomer retrieves all customers
func getCustomers(db *gorm.DB, c *fiber.Ctx) error {
	var customer []Customer
	db.Find(&customer)
	return c.JSON(customer)
}

// getCustomer retrieves a customer by id
func getCustomer(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var customer Customer
	db.First(&customer, id)
	return c.JSON(customer)
}

// updateCustomer updates a Customer by id
func updateCustomer(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	customer := new(Customer)
	db.First(&customer, id)
	if err := c.BodyParser(customer); err != nil {
		return err
	}
	db.Save(&customer)
	return c.JSON(customer)
}

// deleteCustomer deletes a Customer by id
func deleteCustomer(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	db.Delete(&Customer{}, id)
	return c.SendString("Customer successfully deleted")
}
