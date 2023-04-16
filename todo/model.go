package todo

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Todo represents a todo item.
type Todo struct {
	ID          uuid.UUID `json:"id"`
	Text        string    `json:"text"`
	Priority    Priority  `json:"priority"`
	Completed   bool      `json:"completed"`
	TimeCreated time.Time `json:"time_created"`
	TimeUpdated time.Time `json:"time_updated"`
}

// Priority is an enum that represents the different priorities a todo item can
// have.
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// TodoCreateParams are what we require from clients to create a todo item.
type TodoCreateParams struct {
	Text     string   `json:"text"`
	Priority Priority `json:"priority"`
}

// Validate validates the TodoCreateOptions.
func (t TodoCreateParams) Validate() error {
	errs := make([]error, 0)

	if t.Text == "" {
		errs = append(errs, errors.New("missing required field text"))
	}

	if t.Priority != PriorityHigh && t.Priority != PriorityMedium && t.Priority != PriorityLow {
		errs = append(errs, fmt.Errorf(
			"invalid priority %q: must be one of [%v, %v, %v]",
			string(t.Priority),
			PriorityLow,
			PriorityMedium,
			PriorityHigh,
		))
	}

	return NewValidationError(errors.Join(errs...))
}

// TodoUpdateParams represents the information that clients can modify for a
// todo item. Pointers are used to determine whether or not a field was
// provided by the client.
type TodoUpdateParams struct {
	Text      *string   `json:"text"`
	Priority  *Priority `json:"priority"`
	Completed *bool     `json:"completed"`
}

// Validate validates the TodoUpdateOptions.
func (t TodoUpdateParams) Validate() error {
	errs := make([]error, 0)

	if t.Text != nil && *t.Text == "" {
		errs = append(errs, errors.New("missing required field text"))
	}

	if t.Priority != nil {
		if *t.Priority != PriorityHigh && *t.Priority != PriorityMedium && *t.Priority != PriorityLow {
			errs = append(errs, fmt.Errorf(
				"invalid priority %q: must be one of [%v, %v, %v]",
				string(*t.Priority),
				PriorityLow,
				PriorityMedium,
				PriorityHigh,
			))
		}
	}

	return NewValidationError(errors.Join(errs...))
}
