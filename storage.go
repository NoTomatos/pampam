package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	ErrNotFound       = errors.New("не найдено")
	ErrDuplicateEmail = errors.New("email уже существует")
)

type Storage interface {
	// tasks
	GetAllTasks(ctx context.Context) ([]Task, error)
	GetTaskByID(ctx context.Context, id int) (*Task, error)
	CreateTask(ctx context.Context, task *Task) error
	UpdateTask(ctx context.Context, task *Task) error
	DeleteTask(ctx context.Context, id int) error

	// users
	GetAllUsers(ctx context.Context) ([]User, error)
	GetUserByID(ctx context.Context, id int) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int) error

	Close() error
}

// in memory
type MemoryStorage struct {
	tasks     []Task
	users     []User
	idCounter int
	idUser    int
	emails    []Email
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tasks:  []Task{},
		users:  []User{},
		emails: []Email{},
	}
}

func (m *MemoryStorage) Close() error {
	return nil
}

func (m *MemoryStorage) GetAllTasks(ctx context.Context) ([]Task, error) {
	return m.tasks, nil
}

func (m *MemoryStorage) GetTaskByID(ctx context.Context, id int) (*Task, error) {
	for _, task := range m.tasks {
		if task.ID == id {
			return &task, nil
		}
	}
	return nil, ErrNotFound
}

func (m *MemoryStorage) CreateTask(ctx context.Context, task *Task) error {
	m.idCounter++
	task.ID = m.idCounter
	m.tasks = append(m.tasks, *task)
	return nil
}

func (m *MemoryStorage) UpdateTask(ctx context.Context, task *Task) error {
	for i, t := range m.tasks {
		if t.ID == task.ID {
			m.tasks[i] = *task
			return nil
		}
	}
	return ErrNotFound
}

func (m *MemoryStorage) DeleteTask(ctx context.Context, id int) error {
	for i, task := range m.tasks {
		if task.ID == id {
			m.tasks = append(m.tasks[:i], m.tasks[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

func (m *MemoryStorage) CreateUser(ctx context.Context, user *User) error {
	for _, e := range m.emails {
		if e == user.Email {
			return ErrDuplicateEmail
		}
	}

	if err := user.HashPassword(); err != nil {
		return err
	}

	m.idUser++
	user.ID = m.idUser
	m.users = append(m.users, *user)
	m.emails = append(m.emails, user.Email)

	return nil
}

func (m *MemoryStorage) GetAllUsers(ctx context.Context) ([]User, error) {
	var safeUsers []User
	for _, u := range m.users {
		safeUser := u
		safeUser.Hash = ""
		safeUsers = append(safeUsers, safeUser)
	}
	return safeUsers, nil
}

func (m *MemoryStorage) GetUserByID(ctx context.Context, id int) (*User, error) {
	for _, user := range m.users {
		if user.ID == id {
			safeUser := user
			safeUser.Hash = ""
			return &safeUser, nil
		}
	}
	return nil, ErrNotFound
}

func (m *MemoryStorage) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	for _, user := range m.users {
		if string(user.Email) == email {
			return &user, nil
		}
	}
	return nil, ErrNotFound
}

func (m *MemoryStorage) UpdateUser(ctx context.Context, user *User) error {
	for i, u := range m.users {
		if u.ID == user.ID {
			if u.Email != user.Email {
				for _, e := range m.emails {
					if e == user.Email {
						return ErrDuplicateEmail
					}
				}
			}
			if user.Password != "" {
				if err := user.HashPassword(); err != nil {
					return err
				}
			} else {
				user.Hash = u.Hash
			}

			updatedUser := *user
			if user.Password == "" {
				updatedUser.Hash = u.Hash
			}
			m.users[i] = updatedUser
			if u.Email != user.Email {
				for j, e := range m.emails {
					if e == u.Email {
						m.emails[j] = user.Email
						break
					}
				}
			}
			return nil
		}
	}
	return ErrNotFound
}

func (m *MemoryStorage) DeleteUser(ctx context.Context, id int) error {
	for i, user := range m.users {
		if user.ID == id {
			for j, e := range m.emails {
				if e == user.Email {
					m.emails = append(m.emails[:j], m.emails[j+1:]...)
					break
				}
			}
			m.users = append(m.users[:i], m.users[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

// postgres
type PostgresStorage struct {
	db *sqlx.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnvAsInt("DB_PORT", 5432)
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "taskmanager")

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения: %w", err)
	}

	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) Close() error {
	return s.db.Close()
}

func (s *PostgresStorage) GetAllTasks(ctx context.Context) ([]Task, error) {
	var tasks []Task
	query := `SELECT * FROM tasks`
	err := s.db.SelectContext(ctx, &tasks, query)
	return tasks, err
}

func (s *PostgresStorage) GetTaskByID(ctx context.Context, id int) (*Task, error) {
	var task Task
	query := `SELECT * FROM tasks WHERE id = $1`
	err := s.db.GetContext(ctx, &task, query, id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &task, err
}

func (s *PostgresStorage) CreateTask(ctx context.Context, task *Task) error {
	query := `INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3) RETURNING id`
	return s.db.QueryRowContext(ctx, query, task.Title, task.Description, task.Status).Scan(&task.ID)
}

func (s *PostgresStorage) UpdateTask(ctx context.Context, task *Task) error {
	query := `UPDATE tasks SET title=$1, description=$2, status=$3 WHERE id=$4`
	result, err := s.db.ExecContext(ctx, query, task.Title, task.Description, task.Status, task.ID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostgresStorage) DeleteTask(ctx context.Context, id int) error {
	query := `DELETE FROM tasks WHERE id=$1`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostgresStorage) CreateUser(ctx context.Context, user *User) error {
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := s.db.GetContext(ctx, &exists, checkQuery, string(user.Email))
	if err != nil {
		return err
	}

	if exists {
		return ErrDuplicateEmail
	}

	if err := user.HashPassword(); err != nil {
		return err
	}

	query := `INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id`
	return s.db.QueryRowContext(ctx, query, user.Name, string(user.Email), user.Hash).Scan(&user.ID)
}

func (s *PostgresStorage) GetAllUsers(ctx context.Context) ([]User, error) {
	var users []User
	query := `SELECT id, name, email FROM users`
	err := s.db.SelectContext(ctx, &users, query)
	return users, err
}

func (s *PostgresStorage) GetUserByID(ctx context.Context, id int) (*User, error) {
	var user User
	query := `SELECT id, name, email FROM users WHERE id=$1`
	err := s.db.GetContext(ctx, &user, query, id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &user, err
}

func (s *PostgresStorage) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query := `SELECT * FROM users WHERE email=$1`
	err := s.db.GetContext(ctx, &user, query, email)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &user, err
}

func (s *PostgresStorage) UpdateUser(ctx context.Context, user *User) error {
	if user.Email != "" {
		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND id != $2)`
		err := s.db.GetContext(ctx, &exists, checkQuery, string(user.Email), user.ID)
		if err != nil {
			return err
		}

		if exists {
			return ErrDuplicateEmail
		}
	}

	if user.Password != "" {
		if err := user.HashPassword(); err != nil {
			return err
		}
	}

	query := `UPDATE users SET name=$1, email=$2, password_hash=$3 WHERE id=$4`
	result, err := s.db.ExecContext(ctx, query, user.Name, string(user.Email), user.Hash, user.ID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostgresStorage) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id=$1`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// others
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if num, err := strconv.Atoi(value); err == nil {
			return num
		}
	}
	return defaultValue
}

func NewStorage() Storage {
	if getEnv("USE_MEMORY", "false") == "true" {
		fmt.Println("Использую in-memory хранилище")
		return NewMemoryStorage()
	}

	pg, err := NewPostgresStorage()
	if err != nil {
		fmt.Printf("PostgreSQL недоступен: %v\n", err)
		return NewMemoryStorage()
	}

	fmt.Println("Использую PostgreSQL хранилище")
	return pg
}

func initTables(storage Storage) error {
	pg, ok := storage.(*PostgresStorage)
	if !ok {
		return nil
	}

	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash TEXT NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		status VARCHAR(50) NOT NULL
	);
	`

	_, err := pg.db.Exec(query)
	return err
}
