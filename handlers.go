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
	ctx := c.Request.Context()
	tasks, err := storage.GetAllTasks(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении задач",
		})
		return
	}
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

	ctx := c.Request.Context()
	err := storage.CreateTask(ctx, &task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при создании задачи",
		})
		return
	}

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

	ctx := c.Request.Context()
	task, err := storage.GetTaskByID(ctx, id)
	if err != nil {
		if err == ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "задача не найдена",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "ошибка при получении задачи",
			})
		}
		return
	}

	c.JSON(http.StatusOK, task)
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

	updTask.ID = id

	if updTask.Status != "" && !updTask.Status.IsStatusValid() {
		ctx := c.Request.Context()
		currentTask, err := storage.GetTaskByID(ctx, id)
		if err == nil && currentTask != nil {
			updTask.Status = currentTask.Status
		} else {
			updTask.Status = StatusNew
		}
	}

	ctx := c.Request.Context()
	err = storage.UpdateTask(ctx, &updTask)
	if err != nil {
		if err == ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "задача не найдена",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "ошибка при обновлении задачи",
			})
		}
		return
	}

	c.JSON(http.StatusOK, updTask)
}

func DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "неверный ID",
		})
		return
	}

	ctx := c.Request.Context()
	err = storage.DeleteTask(ctx, id)
	if err != nil {
		if err == ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "задача не найдена",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "ошибка при удалении задачи",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "задача удалена",
	})
}

func GetAllUsers(c *gin.Context) {
	ctx := c.Request.Context()
	users, err := storage.GetAllUsers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении пользователей",
		})
		return
	}
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

	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "пароль не задан",
		})
		return
	}

	if !user.Email.IsEmailValid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "неверный формат email",
		})
		return
	}

	ctx := c.Request.Context()
	err := storage.CreateUser(ctx, &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user.Hash = ""
	user.Password = ""

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

	ctx := c.Request.Context()
	user, err := storage.GetUserByID(ctx, id)
	if err != nil {
		if err == ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "пользователь не найден",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "ошибка при получении пользователя",
			})
		}
		return
	}

	c.JSON(http.StatusOK, user)
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

	if !updUser.Email.IsEmailValid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "неверный формат email",
		})
		return
	}

	if err := c.ShouldBindJSON(&updUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "неверный формат JSON",
		})
		return
	}

	updUser.ID = id

	ctx := c.Request.Context()
	err = storage.UpdateUser(ctx, &updUser)
	if err != nil {
		if err == ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "пользователь не найден",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "ошибка при обновлении пользователя",
			})
		}
		return
	}

	updUser.Hash = ""
	updUser.Password = ""

	c.JSON(http.StatusOK, updUser)
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "неверный ID",
		})
		return
	}

	ctx := c.Request.Context()
	err = storage.DeleteUser(ctx, id)
	if err != nil {
		if err == ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "пользователь не найден",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "ошибка при удалении пользователя",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "пользователь удалён",
	})
}
