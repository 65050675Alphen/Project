package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	RoomSize string  // ขนาดห้อง
	PetType  string  // ประเภทสัตว์ที่รองรับ เช่น สุนัขหรือแมว
	Price    float64 // ราคาต่อคืน
	Pets     []Pet   `gorm:"foreignKey:RoomID"`
}

func createRoom(db *gorm.DB, c *fiber.Ctx) error {
	room := new(Room)

	// ดึงข้อมูลจาก request body แล้วทำการ parse เป็น struct Room
	if err := c.BodyParser(room); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// บันทึกข้อมูลห้องใหม่ลงในฐานข้อมูล
	result := db.Create(&room)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": result.Error.Error(),
		})
	}

	// ส่ง response กลับไปพร้อมกับข้อมูลห้องที่ถูกสร้าง
	return c.Status(fiber.StatusCreated).JSON(room)
}

func getRooms(db *gorm.DB, c *fiber.Ctx) error {
	var rooms []Room
	db.Find(&rooms)
	return c.JSON(rooms)
}

func getRoom(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var room Room
	result := db.First(&room, id)
	if result.Error != nil {
		return result.Error
	}
	return c.JSON(room)
}

func updateRoom(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	room := new(Room)
	db.First(&room, id)
	if err := c.BodyParser(room); err != nil {
		return err
	}
	db.Save(&room)
	return c.JSON(room)
}

func deleteRoom(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	db.Delete(&Room{}, id)
	return c.SendString("Room successfully deleted")
}
