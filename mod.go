package main

type Status string

type Email string

const (
	StatusNew        Status = "Новая задача"
	StatusInProgress Status = "Задача в процессе"
	StatusCompleted  Status = "Задача выполнена!"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      Status `json:"status"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    Email  `json:"email"`
	Password string `json:"password"`
}

var (
	tasks     []Task
	idCounter int
	users     []User
	idUser    int
	emails    []Email
)
