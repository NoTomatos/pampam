package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

var storage Storage

func main() {
	storage = NewStorage()

	defer func() {
		if err := storage.Close(); err != nil {
			log.Printf("Ошибка при закрытии хранилища: %v", err)
		}
		fmt.Println("Хранилище закрыто")
	}()

	if err := initTables(storage); err != nil {
		log.Printf("Не удалось создать таблицы: %v", err)
		log.Println("Продолжаю работу без таблиц...")
	} else {
		fmt.Println("Таблицы проверены/созданы")
	}

	r := gin.Default()

	fmt.Println("Регистрирую маршруты для задач...")
	r.GET("/tasks", GetAllTasks)
	r.GET("/tasks/:id", GetTaskByID)
	r.POST("/tasks", CreateTask)
	r.PUT("/tasks/:id", UpdateTask)
	r.DELETE("/tasks/:id", DeleteTask)

	fmt.Println("Регистрирую маршруты для пользователей...")
	r.GET("/users", GetAllUsers)
	r.GET("/users/:id", GetUserByID)
	r.POST("/users", CreateUser)
	r.PUT("/users/:id", UpdateUser)
	r.DELETE("/users/:id", DeleteUser)

	port := getEnv("SERVER_PORT", ":8080")
	fmt.Printf("Сервер слушает на порту %s\n", port)

	if err := r.Run(port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
