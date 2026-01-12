package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/tasks", GetAllTasks)
	r.GET("/tasks/:id", GetTaskByID)
	r.POST("/tasks", CreateTask)
	r.PUT("/tasks/:id", UpdateTask)
	r.DELETE("/tasks/:id", DeleteTask)
	r.GET("/users", GetAllUsers)
	r.GET("/users/:id", GetUserByID)
	r.POST("/users", CreateUser)
	r.PUT("/users/:id", UpdateUser)
	r.DELETE("/users/:id", DeleteUser)
	r.Run(":8080")
}
