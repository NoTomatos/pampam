package main

type Status string

type Email string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
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
	Hash     string `json:"-" db:"password_hash"`
}
