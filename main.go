package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"myapi/models"
)

var db *gorm.DB
var err error

// Função para inicializar a conexão com o banco de dados
func initDatabase() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database! \n", err)
	}

	// Faz o auto-migrate das tabelas User e Address
	db.AutoMigrate(&models.User{}, &models.Address{})
}

func main() {
	// Inicializa o banco de dados
	initDatabase()

	// Cria a instância do Fiber
	app := fiber.New()

	// Rotas CRUD para User
	app.Get("/users", GetUsers)
	app.Get("/users/:id", GetUserByID)
	app.Post("/users", CreateUser)
	app.Put("/users/:id", UpdateUser)
	app.Delete("/users/:id", DeleteUser)

	// Rotas CRUD para Address
	app.Get("/addresses", GetAddresses)
	app.Get("/addresses/:id", GetAddressByID)
	app.Post("/addresses", CreateAddress)
	app.Put("/addresses/:id", UpdateAddress)
	app.Delete("/addresses/:id", DeleteAddress)

	// Inicia o servidor na porta 3000
	log.Fatal(app.Listen(":3000"))
}

// Handlers para User
func GetUsers(c *fiber.Ctx) error {
	var users []models.User
	db.Preload("Address").Find(&users) // Preload carrega a relação com Address
	return c.JSON(users)
}

func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	if err := db.Preload("Address").First(&user, id).Error; err != nil {
		return c.Status(404).SendString("User not found")
	}
	return c.JSON(user)
}

func CreateUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := db.Create(&user).Error; err != nil {
		return c.Status(500).SendString("Could not create user")
	}
	return c.JSON(user)
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return c.Status(404).SendString("User not found")
	}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	db.Save(&user)
	return c.JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return c.Status(404).SendString("User not found")
	}
	db.Delete(&user)
	return c.SendString("User deleted")
}

// Handlers para Address
func GetAddresses(c *fiber.Ctx) error {
	var addresses []models.Address
	db.Find(&addresses)
	return c.JSON(addresses)
}

func GetAddressByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var address models.Address
	if err := db.First(&address, id).Error; err != nil {
		return c.Status(404).SendString("Address not found")
	}
	return c.JSON(address)
}

func CreateAddress(c *fiber.Ctx) error {
	address := new(models.Address)
	if err := c.BodyParser(address); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := db.Create(&address).Error; err != nil {
		return c.Status(500).SendString("Could not create address")
	}
	return c.JSON(address)
}

func UpdateAddress(c *fiber.Ctx) error {
	id := c.Params("id")
	var address models.Address
	if err := db.First(&address, id).Error; err != nil {
		return c.Status(404).SendString("Address not found")
	}
	if err := c.BodyParser(&address); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	db.Save(&address)
	return c.JSON(address)
}

func DeleteAddress(c *fiber.Ctx) error {
	id := c.Params("id")
	var address models.Address
	if err := db.First(&address, id).Error; err != nil {
		return c.Status(404).SendString("Address not found")
	}
	db.Delete(&address)
	return c.SendString("Address deleted")
}
