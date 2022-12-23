package todo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("todo not found")
)

// Storer represents the behavior this package needs to manage todo items.
type Storer interface {
	Query(ctx context.Context) ([]Todo, error)
	QueryByID(ctx context.Context, id uuid.UUID) (Todo, error)
	Create(ctx context.Context, todo Todo) error
	Update(ctx context.Context, todo Todo) error
	Delete(ctx context.Context, todo Todo) error
}

// Core exposes the APIs needed to interface with todo items.
type Core struct {
	storer Storer
}

// NewCore is a constructor for a Core.
func NewCore(storer Storer) *Core {
	return &Core{
		storer: storer,
	}
}

// Query retrieves all todo items.
func (s *Core) Query(ctx context.Context) ([]Todo, error) {
	todos, err := s.storer.Query(ctx)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return todos, nil
}

// QueryByID retrieves a todo item by its ID.
func (s *Core) QueryByID(ctx context.Context, id uuid.UUID) (Todo, error) {
	t, err := s.storer.QueryByID(ctx, id)
	if err != nil {
		return Todo{}, fmt.Errorf("query by id: %w", err)
	}

	return t, nil
}

// Create adds a todo item into the store.
func (s *Core) Create(ctx context.Context, params TodoCreateParams) (Todo, error) {
	if err := params.Validate(); err != nil {
		return Todo{}, fmt.Errorf("validate: %w", err)
	}

	now := time.Now()

	todo := Todo{
		ID:          uuid.New(),
		Text:        params.Text,
		Priority:    params.Priority,
		Completed:   false,
		TimeCreated: now,
		TimeUpdated: now,
	}

	if err := s.storer.Create(ctx, todo); err != nil {
		return Todo{}, fmt.Errorf("create: %w", err)
	}

	return todo, nil
}

// Update modifies an existing todo item.
func (s *Core) Update(ctx context.Context, todo Todo, params TodoUpdateParams) (Todo, error) {
	if err := params.Validate(); err != nil {
		return Todo{}, fmt.Errorf("validate: %w", err)
	}

	if params.Text != nil {
		todo.Text = *params.Text
	}
	if params.Priority != nil {
		todo.Priority = *params.Priority
	}
	if params.Completed != nil {
		todo.Completed = *params.Completed
	}

	todo.TimeUpdated = time.Now()

	if err := s.storer.Update(ctx, todo); err != nil {
		return Todo{}, fmt.Errorf("update: %w", err)
	}

	return todo, nil
}

// Delete deletes the specified todo item.
func (s *Core) Delete(ctx context.Context, todo Todo) error {
	if err := s.storer.Delete(ctx, todo); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
