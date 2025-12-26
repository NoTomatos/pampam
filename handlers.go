package main

type TaskStore struct {
	tasks     []Task
	idCounter int
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks:     make([]Task, 0),
		idCounter: 0,
	}
}

func (t *TaskStore) GetAllTasks() []Task {
	return t.tasks
}

func (t *TaskStore) GetTaskByID(id int) (*Task, bool) {
	for _, task := range t.tasks {
		if task.ID == id {
			return &task, true
		}
	}
	return nil, false
}

func (t *TaskStore) CreateTask(task Task) Task {
	t.idCounter++
	task.ID = t.idCounter
	if task.Status != "" {
		task.Status = StatusNew
	}
	t.tasks = append(t.tasks, task)
	return task
}
