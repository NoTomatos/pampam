package main

type Status string

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
