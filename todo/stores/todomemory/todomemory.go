package todomemory

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/sudomateo/todo/todo"
)

// Store exposes the APIs needed to interface with todo items in memory.
type Store struct {
	data  []todo.Todo
	mutex sync.RWMutex
}

// NewStore is a constructor for a Store.
func NewStore() *Store {
	return &Store{
		data: make([]todo.Todo, 0),
	}
}

// Query retrieves all the todo items from memory.
func (d *Store) Query(ctx context.Context) ([]todo.Todo, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.data, nil
}

// QueryByID retrieves a todo item from memory.
func (d *Store) QueryByID(ctx context.Context, id uuid.UUID) (todo.Todo, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	for i := range d.data {
		if d.data[i].ID == id {
			return d.data[i], nil
		}
	}

	return todo.Todo{}, todo.ErrNotFound
}

// Create adds a todo item to memory.
func (d *Store) Create(ctx context.Context, td todo.Todo) error {
	d.mutex.Lock()

	d.data = append(d.data, todo.Todo{
		ID:          td.ID,
		Text:        td.Text,
		Priority:    td.Priority,
		Completed:   td.Completed,
		TimeCreated: td.TimeCreated,
		TimeUpdated: td.TimeUpdated,
	})

	d.mutex.Unlock()

	return nil
}

// Update modifies an existing todo item in memory.
func (d *Store) Update(ctx context.Context, td todo.Todo) error {
	d.mutex.Lock()

	for i := range d.data {
		if d.data[i].ID == td.ID {
			d.data[i].Text = td.Text
			d.data[i].Priority = td.Priority
			d.data[i].Completed = td.Completed
			d.data[i].TimeUpdated = td.TimeUpdated
		}
	}

	d.mutex.Unlock()

	return nil
}

// Delete deletes a todo item from memory.
func (d *Store) Delete(ctx context.Context, td todo.Todo) error {
	d.mutex.Lock()

	for i := range d.data {
		if d.data[i].ID == td.ID {
			d.data = append(d.data[:i], d.data[i+1:]...)
			break
		}
	}

	d.mutex.Unlock()

	return nil
}
