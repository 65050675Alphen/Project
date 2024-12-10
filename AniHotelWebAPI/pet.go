package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Pet Table
type Pet struct {
	gorm.Model
	Petname     string `json:"name"`
	PetType     string `json:"type"`
	Description string `json:"description"`
	CustomerID  uint   `json:"customer_id"`
}

// createPetWithCustomerID creates a pet for a specific customer
func createPetWithCustomerID(db *gorm.DB, c *fiber.Ctx) error {
	customerID := c.Params("customerID") // รับ CustomerID จาก URL path
	pet := new(Pet)

	// Parse request body into Pet struct
	if err := c.BodyParser(pet); err != nil {
		return err
	}

	// ตรวจสอบว่า CustomerID นั้นมีอยู่จริงหรือไม่
	var customer Customer
	if err := db.First(&customer, customerID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Customer not found")
	}

	// Set the CustomerID in the pet struct
	pet.CustomerID = customer.ID

	// สร้าง record ของสัตว์เลี้ยงในฐานข้อมูล
	if err := db.Create(&pet).Error; err != nil {
		return err
	}

	return c.JSON(pet)
}

// getPets retrieves all pets for a specific customer
func getPets(db *gorm.DB, c *fiber.Ctx) error {
	customerID := c.Params("customerID") // รับ CustomerID จาก URL path
	var pets []Pet
	if err := db.Where("customer_id = ?", customerID).Find(&pets).Error; err != nil {
		return err
	}
	return c.JSON(pets)
}

// getPet retrieves a pet by id for a specific customer
func getPet(db *gorm.DB, c *fiber.Ctx) error {
	customerID := c.Params("customerID") // รับ CustomerID จาก URL path
	petID := c.Params("id")
	var pet Pet
	if err := db.First(&pet, "id = ? AND customer_id = ?", petID, customerID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Pet not found for this customer")
	}
	return c.JSON(pet)
}

// updatePet updates a pet by id for a specific customer
func updatePet(db *gorm.DB, c *fiber.Ctx) error {
	customerID := c.Params("customerID") // รับ CustomerID จาก URL path
	petID := c.Params("id")
	var pet Pet

	// ตรวจสอบว่า CustomerID นั้นมีอยู่จริงหรือไม่
	if err := db.First(&pet, "id = ? AND customer_id = ?", petID, customerID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Pet not found for this customer")
	}

	// Update Pet details
	if err := c.BodyParser(&pet); err != nil {
		return err
	}
	db.Save(&pet)

	return c.JSON(pet)
}

// deletePet deletes a pet by id for a specific customer
func deletePet(db *gorm.DB, c *fiber.Ctx) error {
	customerID := c.Params("customerID") // รับ CustomerID จาก URL path
	petID := c.Params("id")

	// ตรวจสอบว่า Pet นั้นมีอยู่และเป็นของ Customer นี้หรือไม่
	if err := db.Delete(&Pet{}, "id = ? AND customer_id = ?", petID, customerID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Pet not found or doesn't belong to this customer")
	}

	return c.SendString("Pet successfully deleted")
}
