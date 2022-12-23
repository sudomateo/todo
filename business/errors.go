package business

import "errors"

type ValidationError struct {
	Err error
}

func NewValidationError(err error) error {
	return &ValidationError{err}
}

func (ve *ValidationError) Error() string {
	return ve.Err.Error()
}

func IsValidationError(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}

func GetValidationError(err error) *ValidationError {
	var ve *ValidationError
	if !errors.As(err, &ve) {
		return nil
	}
	return ve
}
