package tododb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/sudomateo/todo/todo"
)

// Store exposes the APIs needed to interface with todo items in the database.
type Store struct {
	db *sql.DB
}

// NewStore is a constructor for a Store.
func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// Query retrieves all the todo items from the database.
func (d *Store) Query(ctx context.Context) ([]todo.Todo, error) {
	query := `SELECT * FROM todos ORDER BY time_created`

	todos := make([]todo.Todo, 0)

	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, todo.ErrNotFound
		}
		return nil, fmt.Errorf("db: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		td := new(todo.Todo)
		if err := rows.Scan(
			&td.ID,
			&td.Text,
			&td.Priority,
			&td.Completed,
			&td.TimeCreated,
			&td.TimeUpdated,
		); err != nil {
			return nil, err
		}

		todos = append(todos, *td)
	}

	return todos, nil
}

// QueryByID retrieves a todo item from the database.
func (d *Store) QueryByID(ctx context.Context, id uuid.UUID) (todo.Todo, error) {
	const query = `SELECT * FROM todos WHERE id = $1 LIMIT 1`

	var t todo.Todo

	if err := d.db.QueryRowContext(ctx, query, id.String()).Scan(
		&t.ID,
		&t.Text,
		&t.Priority,
		&t.Completed,
		&t.TimeCreated,
		&t.TimeUpdated,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return todo.Todo{}, todo.ErrNotFound
		}
		return todo.Todo{}, fmt.Errorf("db: %w", err)
	}

	return t, nil
}

// Create adds a todo item to the database.
func (d *Store) Create(ctx context.Context, td todo.Todo) error {
	const query = `
	INSERT INTO todos
	  (id, text, priority, completed, time_created, time_updated)
	VALUES
	  ($1, $2, $3, $4, $5, $6)`

	if _, err := d.db.ExecContext(ctx, query,
		td.ID,
		td.Text,
		td.Priority,
		td.Completed,
		td.TimeCreated,
		td.TimeUpdated,
	); err != nil {
		return fmt.Errorf("db: %w", err)
	}

	return nil
}

// Update modifies an existing todo item in the database.
func (d *Store) Update(ctx context.Context, td todo.Todo) error {
	const query = `
	UPDATE
	  todos
	SET
		text = $1,
		priority = $2,
		completed = $3,
		time_updated = $4
	WHERE
	  id = $5`

	if _, err := d.db.ExecContext(ctx, query,
		td.Text,
		td.Priority,
		td.Completed,
		td.TimeUpdated,
		td.ID,
	); err != nil {
		return fmt.Errorf("db: %w", err)
	}

	return nil
}

// Delete deletes a todo item from the database.
func (d *Store) Delete(ctx context.Context, td todo.Todo) error {
	const query = `
	DELETE FROM
	  todos
	WHERE
	  id = $1`

	if _, err := d.db.ExecContext(ctx, query,
		td.ID,
	); err != nil {
		return fmt.Errorf("db: %w", err)
	}

	return nil
}
