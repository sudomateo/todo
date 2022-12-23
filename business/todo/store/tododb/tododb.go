package tododb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"todo-api/business/todo"

	"github.com/google/uuid"
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
func (d *Store) Query(ctx context.Context, filters []string) ([]todo.Todo, error) {
	query := `SELECT * FROM todos %s ORDER BY id OFFSET 0 ROWS FETCH NEXT 20 ROWS ONLY`

	if len(filters) > 0 {
		query = fmt.Sprintf(query, "WHERE "+strings.Join(filters, " AND "))
	} else {
		query = fmt.Sprintf(query, "")
	}

	todos := make([]todo.Todo, 0, 20)

	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, todo.ErrNotFound
		}
		return nil, fmt.Errorf("querying todo: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		dbTodo := new(dbTodo)
		if err := rows.Scan(
			&dbTodo.ID,
			&dbTodo.Text,
			&dbTodo.Priority,
			&dbTodo.Completed,
			&dbTodo.TimeCreated,
			&dbTodo.TimeUpdated,
		); err != nil {
			return nil, err
		}

		todos = append(todos, toTodo(*dbTodo))
	}

	return todos, nil
}

// Create adds a todo item to the database.
func (d *Store) Create(ctx context.Context, td todo.Todo) error {
	const query = `
	INSERT INTO todos
	  (id, text, priority, completed, time_created, time_updated)
	VALUES
	  ($1, $2, $3, $4, $5, $6)`

	dbTodo := toDBTodo(td)

	if _, err := d.db.ExecContext(ctx, query,
		dbTodo.ID,
		dbTodo.Text,
		dbTodo.Priority,
		dbTodo.Completed,
		dbTodo.TimeCreated,
		dbTodo.TimeUpdated,
	); err != nil {
		return fmt.Errorf("inserting todo: %w", err)
	}

	return nil
}

// QueryByID retrieves a todo item from the database.
func (d *Store) QueryByID(ctx context.Context, todoID uuid.UUID) (todo.Todo, error) {
	const query = `SELECT * FROM todos WHERE id = $1 LIMIT 1`

	var dbTodo dbTodo

	if err := d.db.QueryRowContext(ctx, query, todoID.String()).Scan(
		&dbTodo.ID,
		&dbTodo.Text,
		&dbTodo.Priority,
		&dbTodo.Completed,
		&dbTodo.TimeCreated,
		&dbTodo.TimeUpdated,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return todo.Todo{}, todo.ErrNotFound
		}
		return todo.Todo{}, fmt.Errorf("querying todo [%v] %w", todoID, err)
	}

	return toTodo(dbTodo), nil
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

	dbTodo := toDBTodo(td)

	if _, err := d.db.ExecContext(ctx, query,
		dbTodo.Text,
		dbTodo.Priority,
		dbTodo.Completed,
		dbTodo.TimeUpdated,
		dbTodo.ID,
	); err != nil {
		return fmt.Errorf("updating todo [%v]: %w", dbTodo.ID, err)
	}

	return nil
}

// Delete delets a todo item from the database.
func (d *Store) Delete(ctx context.Context, td todo.Todo) error {
	const query = `
	DELETE FROM
	  todos
	WHERE
	  id = $1`

	dbTodo := toDBTodo(td)

	if _, err := d.db.ExecContext(ctx, query,
		dbTodo.ID,
	); err != nil {
		return fmt.Errorf("deleting todo [%v]: %w", dbTodo.ID, err)
	}

	return nil
}
