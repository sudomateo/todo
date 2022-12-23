package todo_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/sudomateo/todo/todo"
	"github.com/sudomateo/todo/todo/stores/todomemory"
)

func TestTodo(t *testing.T) {
	todoCore := todo.NewCore(todomemory.NewStore())

	todos, err := todoCore.Query(context.Background())
	if err != nil {
		t.Fatalf("query: expected nil error, got %v", err)
	}

	if diff := cmp.Diff([]todo.Todo{}, todos); diff != "" {
		t.Fatalf("query: did not return empty array of todos: %v", diff)
	}

	td, err := todoCore.Create(context.Background(), todo.TodoCreateParams{
		Text:     "foo",
		Priority: todo.PriorityHigh,
	})
	if err != nil {
		t.Fatalf("create: expected nil error, got %v", err)
	}

	tdQuery, err := todoCore.QueryByID(context.Background(), td.ID)
	if err != nil {
		t.Fatalf("query by id: expected nil error, got %v", err)
	}

	if diff := cmp.Diff(td, tdQuery); diff != "" {
		t.Fatalf("compare: %v", diff)
	}

	todos, err = todoCore.Query(context.Background())
	if err != nil {
		t.Fatalf("query: expected nil error, got %v", err)
	}
	if len(todos) != 1 {
		t.Fatalf("query: expected 1 todo, got %v", len(todos))
	}

	if err := todoCore.Delete(context.Background(), td); err != nil {
		t.Fatalf("delete: expected nil error, got %v", err)
	}

	if err := todoCore.Delete(context.Background(), td); err != nil {
		t.Fatalf("delete: expected nil error, got %v", err)
	}
}
