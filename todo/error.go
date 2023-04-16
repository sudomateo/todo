package todo

// ValidationError represents an error with user input validation.
type ValidationError struct {
	Err error
}

// NewValidationError returns a new validation error.
func NewValidationError(err error) error {
	return ValidationError{
		Err: err,
	}
}

// Error implments the error interface for ValidationError.
func (v ValidationError) Error() string {
	return v.Err.Error()
}
