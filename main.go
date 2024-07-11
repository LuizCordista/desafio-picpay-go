package main

import (
	"desafio-picpay/handlers"
	"desafio-picpay/models"
	"desafio-picpay/repositories"
	"desafio-picpay/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

func main() {
	db, err := gorm.Open("postgres", "user=postgres dbname=picpay password=123 sslmode=disable")
	if err != nil {
		fmt.Println(err)
		panic("Failed to connect to database")
	}
	defer db.Close()

	if err := db.AutoMigrate(&models.User{}).Error; err != nil {
		fmt.Println("AutoMigrate error:", err)
		panic("Failed to auto migrate database")
	}

	userRepository := repositories.NewGormUserRepository(db)
	userService := services.NewUserService(*userRepository)
	userHandler := handlers.NewUserHandler(userService)

	r := gin.Default()

	r.POST("/users", userHandler.CreateUser)
	r.POST("/transfer", userHandler.Transfer)

	r.Run()
}
