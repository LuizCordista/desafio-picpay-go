package models

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const secret = "secret"

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key;unique;not null" json:"id"`
	Name      string    `json:"name" binding:"required"`
	IsShop    bool      `json:"is_shop"`
	CPF       string    `gorm:"unique;not null" json:"cpf" binding:"required"`
	Email     string    `gorm:"unique;not null" json:"email" binding:"required" `
	Password  string    `json:"password" binding:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Currency  float64    `gorm:"default:0" json:"currency"`
}

func EncryptPassword(password string) (string, error) {
	password = password + secret
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(userPassword string, providedPassword string) bool {
	providedPassword = providedPassword + secret
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	return err == nil
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {

	password := user.Password
	encryptedPassword, err := EncryptPassword(password)
	if err != nil {
		fmt.Println("Erro ao encriptografar a senha:", err)
		return err
	}
	fmt.Println("Senha encriptografada:", encryptedPassword)
	user.Password = encryptedPassword
	scope.SetColumn("Password", encryptedPassword)
	return nil
}

type TransferRequest struct {
	Value float64 `json:"value"`
	Payer string  `json:"payer"`
	Payee string  `json:"payee"`
}
