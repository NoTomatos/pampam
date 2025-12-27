package main

import "github.com/gin-gonic/gin"

func main() {
	tasks = make([]Task, 0)
	idCounter = 0
	r := gin.Default()
	r.GET("/tasks", GetAllTasks)
	r.GET("/tasks/:id", GetTaskByID)
	r.POST("/tasks", CreateTask)
	r.PUT("/tasks/:id", UpdateTask)
	r.DELETE("/tasks/:id", DeleteTask)
	r.Run(":8080")
}
