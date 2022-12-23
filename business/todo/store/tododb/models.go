package tododb

import (
	"time"
	"todo-api/business/todo"

	"github.com/google/uuid"
)

// dbTodo represents a todo item as seen by the database.
type dbTodo struct {
	ID          string
	Text        string
	Priority    string
	Completed   bool
	TimeCreated time.Time
	TimeUpdated time.Time
}

func toDBTodo(td todo.Todo) dbTodo {
	dbTodo := dbTodo{
		ID:          td.ID.String(),
		Text:        td.Text,
		Priority:    string(td.Priority),
		Completed:   td.Completed,
		TimeCreated: td.TimeCreated.UTC(),
		TimeUpdated: td.TimeUpdated.UTC(),
	}

	return dbTodo
}

func toTodo(dbTodo dbTodo) todo.Todo {
	td := todo.Todo{
		ID:          uuid.MustParse(dbTodo.ID),
		Text:        dbTodo.Text,
		Priority:    todo.Priority(dbTodo.Priority),
		Completed:   dbTodo.Completed,
		TimeCreated: dbTodo.TimeCreated.In(time.Local),
		TimeUpdated: dbTodo.TimeUpdated.In(time.Local),
	}

	return td
}
