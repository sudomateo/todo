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

// TodoCreateOptions are what we require from clients to create a todo item.
type TodoCreateOptions struct {
	Text     string   `json:"text"`
	Priority Priority `json:"priority"`
}

// Validate validates the TodoCreateOptions.
func (t TodoCreateOptions) Validate() error {
	if t.Text == "" {
		return errors.New("missing required field text")
	}

	if t.Priority == "" {
		return errors.New("missing required field priority")
	}

	if t.Priority != PriorityHigh && t.Priority != PriorityMedium && t.Priority != PriorityLow {
		return fmt.Errorf(
			"invalid priority %q: must be one of [%v, %v, %v]",
			string(t.Priority),
			PriorityLow,
			PriorityMedium,
			PriorityHigh,
		)
	}

	return nil
}

// TodoUpdateOptions represents the information that clients can modify for a
// todo item. Pointers are used to determine whether or not a field was
// provided by the client.
type TodoUpdateOptions struct {
	Text      *string   `json:"text"`
	Priority  *Priority `json:"priority"`
	Completed *bool     `json:"completed"`
}

// Validate validates the TodoUpdateOptions.
func (t TodoUpdateOptions) Validate() error {
	if t.Text != nil && *t.Text == "" {
		return errors.New("missing required field text")
	}

	if t.Priority != nil && *t.Priority == "" {
		return errors.New("missing required field priority")
	}

	if t.Priority != nil {
		if *t.Priority != PriorityHigh && *t.Priority != PriorityMedium && *t.Priority != PriorityLow {
			return fmt.Errorf(
				"invalid priority %q: must be one of [%v, %v, %v]",
				string(*t.Priority),
				PriorityLow,
				PriorityMedium,
				PriorityHigh,
			)
		}
	}

	return nil
}
