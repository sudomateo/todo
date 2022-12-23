package todo

import (
	"context"
	"errors"
	"fmt"
	"time"
	"todo-api/business"

	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("todo not found")
)

// Storer represents the behavior this package needs to manage todo items.
type Storer interface {
	Query(context.Context, []string) ([]Todo, error)
	Create(context.Context, Todo) error
	QueryByID(context.Context, uuid.UUID) (Todo, error)
	Update(context.Context, Todo) error
	Delete(context.Context, Todo) error
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
func (s *Core) Query(ctx context.Context, filters []string) ([]Todo, error) {
	todos, err := s.storer.Query(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return todos, nil
}

// Create adds a todo item into the store.
func (s *Core) Create(ctx context.Context, opts TodoCreateOptions) (Todo, error) {
	if err := opts.Validate(); err != nil {
		return Todo{}, business.NewValidationError(fmt.Errorf("validate: %w", err))
	}

	now := time.Now()

	todo := Todo{
		ID:          uuid.New(),
		Text:        opts.Text,
		Priority:    opts.Priority,
		Completed:   false,
		TimeCreated: now,
		TimeUpdated: now,
	}

	if err := s.storer.Create(ctx, todo); err != nil {
		return Todo{}, fmt.Errorf("create: %w", err)
	}

	return todo, nil
}

// QueryByID retrieves a todo item by its ID.
func (s *Core) QueryByID(ctx context.Context, id uuid.UUID) (Todo, error) {
	t, err := s.storer.QueryByID(ctx, id)
	if err != nil {
		return Todo{}, fmt.Errorf("query: %w", err)
	}

	return t, nil
}

// Update modifies an existing todo item.
func (s *Core) Update(ctx context.Context, todo Todo, opts TodoUpdateOptions) (Todo, error) {
	if err := opts.Validate(); err != nil {
		return Todo{}, business.NewValidationError(fmt.Errorf("validate: %w", err))
	}

	if opts.Text != nil {
		todo.Text = *opts.Text
	}
	if opts.Priority != nil {
		todo.Priority = *opts.Priority
	}
	if opts.Completed != nil {
		todo.Completed = *opts.Completed
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
