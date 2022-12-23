package todohandlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"todo-api/business"
	"todo-api/business/todo"
	v1 "todo-api/business/web/v1"
	"todo-api/web"

	"github.com/google/uuid"
)

var (
	ErrInvalidID = errors.New("invalid todo id")
)

type Handler struct {
	Todo *todo.Core
}

func (h Handler) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	queryParams := r.URL.Query()

	filters := make([]string, 0)

	completedParam := queryParams.Get("completed")
	if completedParam != "" {
		completed, err := strconv.ParseBool(completedParam)
		if err != nil {
			return business.NewValidationError(fmt.Errorf("invalid query parameter completed=%v: must be one of [true, false]", completedParam))
		}

		filters = append(filters, fmt.Sprintf("completed = %t", completed))
	}

	priorityParam := queryParams.Get("priority")
	if priorityParam != "" {
		priority := todo.Priority(priorityParam)
		if priority != todo.PriorityHigh && priority != todo.PriorityMedium && priority != todo.PriorityLow {
			return business.NewValidationError(fmt.Errorf(
				"invalid query parameter priority=%v: must be one of [%v, %v, %v]",
				priority,
				todo.PriorityLow,
				todo.PriorityMedium,
				todo.PriorityHigh,
			))
		}

		filters = append(filters, fmt.Sprintf("priority LIKE '%%%s%%'", priority))
	}

	t, err := h.Todo.Query(ctx, filters)
	if err != nil {
		switch {
		case errors.Is(err, todo.ErrNotFound):
			return v1.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying todo: %w", err)
		}
	}

	return web.Respond(ctx, w, t, http.StatusOK)
}

func (h Handler) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var opts todo.TodoCreateOptions
	if err := web.Decode(r, &opts); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	t, err := h.Todo.Create(ctx, opts)
	if err != nil {
		return fmt.Errorf("creating new todo [%v]: %w", opts, err)
	}

	return web.Respond(ctx, w, t, http.StatusCreated)
}

func (h Handler) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	todoID, err := uuid.Parse(web.Param(r, "id"))
	if err != nil {
		return v1.NewRequestError(ErrInvalidID, http.StatusBadRequest)
	}

	t, err := h.Todo.QueryByID(ctx, todoID)
	if err != nil {
		switch {
		case errors.Is(err, todo.ErrNotFound):
			return v1.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying todo [%v]: %w", todoID, err)
		}
	}

	return web.Respond(ctx, w, t, http.StatusOK)
}

func (h Handler) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var opts todo.TodoUpdateOptions
	if err := web.Decode(r, &opts); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	todoID, err := uuid.Parse(web.Param(r, "id"))
	if err != nil {
		return v1.NewRequestError(ErrInvalidID, http.StatusBadRequest)
	}

	t, err := h.Todo.QueryByID(ctx, todoID)
	if err != nil {
		switch {
		case errors.Is(err, todo.ErrNotFound):
			return v1.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying todo [%v]: %w", todoID, err)
		}
	}

	t, err = h.Todo.Update(ctx, t, opts)
	if err != nil {
		return fmt.Errorf("updating todo [%v] [%+v]: %w", todoID, &opts, err)
	}

	return web.Respond(ctx, w, t, http.StatusOK)
}

func (h Handler) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	todoID, err := uuid.Parse(web.Param(r, "id"))
	if err != nil {
		return v1.NewRequestError(ErrInvalidID, http.StatusBadRequest)
	}

	t, err := h.Todo.QueryByID(ctx, todoID)
	if err != nil {
		switch {
		case errors.Is(err, todo.ErrNotFound):
			return web.Respond(ctx, w, nil, http.StatusNoContent)
		default:
			return fmt.Errorf("querying todo [%v]: %w", todoID, err)
		}
	}

	if err = h.Todo.Delete(ctx, t); err != nil {
		return fmt.Errorf("deleting todo [%v]: %w", todoID, err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
