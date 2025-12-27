package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	tasks     []Task
	idCounter int
)

func (s Status) IsValid() bool {
	if s == StatusNew || s == StatusInProgress || s == StatusCompleted {
		return true
	}
	return false
}

func GetAllTasks(c *gin.Context) {
	c.JSON(http.StatusOK, tasks)
}

func CreateTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "неверный формат JSON",
		})
		return
	}
	if task.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "пустой заголовок",
		})
		return
	}
	if task.Status == "" || !task.Status.IsValid() {
		task.Status = StatusNew
	}
	idCounter++
	task.ID = idCounter
	tasks = append(tasks, task)
	c.JSON(http.StatusCreated, task)
}

func GetTaskByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "неверный ID",
		})
		return
	}
	for _, task := range tasks {
		if task.ID == id {
			c.JSON(http.StatusOK, task)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": "задача не найдена",
	})
}

func UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "неверный ID",
		})
		return
	}
	var updTask Task
	if err := c.ShouldBindJSON(&updTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "неверный формат JSON",
		})
		return
	}
	for i, task := range tasks {
		if task.ID == id {
			updTask.ID = id
			if updTask.Status == "" || !updTask.Status.IsValid() {
				updTask.Status = task.Status
			}
			tasks[i] = updTask
			c.JSON(http.StatusOK, updTask)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": "задача не найдена",
	})
}

func DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "неверный ID",
		})
		return
	}
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			c.JSON(http.StatusOK, gin.H{
				"message": "задача удалена",
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": "задача не найдена",
	})
}
