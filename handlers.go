package main

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (s Status) IsStatusValid() bool {
	if s == StatusNew || s == StatusInProgress || s == StatusCompleted {
		return true
	}
	return false
}

func (e Email) IsEmailValid() bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(string(e))
}

func (u *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Hash = string(hash)
	u.Password = ""
	return nil
}

func (u *User) CheckPassword(pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Hash), []byte(pw))
	return err == nil
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
	if task.Status == "" || !task.Status.IsStatusValid() {
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
			if updTask.Status == "" || !updTask.Status.IsStatusValid() {
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

func GetAllUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "неверный формат JSON",
		})
		return
	}
	if user.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "имя не задано",
		})
		return
	}
	if user.Email == "" || !user.Email.IsEmailValid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "неверный формат email или такой email уже есть",
		})
		return
	}
	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "пароль не задан",
		})
		return
	}
	idUser++
	emails = append(emails, user.Email)
	user.ID = idUser
	users = append(users, user)
	c.JSON(http.StatusCreated, user)
}

func GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "неверный ID",
		})
		return
	}
	for _, user := range users {
		if user.ID == id {
			c.JSON(http.StatusOK, user)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": "пользователь не найден",
	})
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "неверный ID",
		})
		return
	}
	var updUser User
	if err := c.ShouldBindJSON(&updUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "неверный формат JSON",
		})
		return
	}
	for i, user := range users {
		if user.ID == id {
			updUser.ID = id
			if updUser.Email == "" || !updUser.Email.IsEmailValid() {
				updUser.Email = user.Email
			}
			if user.Name == "" {
				updUser.Name = user.Name
			}
			users[i] = updUser
			c.JSON(http.StatusOK, updUser)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": "пользователь не найден",
	})
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "неверный ID",
		})
		return
	}
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			c.JSON(http.StatusOK, gin.H{
				"message": "пользователь удалён",
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": "пользователь не найден",
	})
}
