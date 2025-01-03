package models

import "gorm.io/gorm"

// Address representa o endereço do usuário
type Address struct {
	gorm.Model
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
}

// User representa um usuário
type User struct {
	gorm.Model
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	AddressID uint    `json:"address_id"` // Chave estrangeira para Address
	Address   Address `json:"address" gorm:"foreignKey:AddressID"`
}
