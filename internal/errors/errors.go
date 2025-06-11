package errors

import "fmt"

// ValidationError represents an error during validation
type ValidationError struct {
	message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s", e.message)
}

// NewValidationError creates a new ValidationError
func NewValidationError(message string) *ValidationError {
	return &ValidationError{message: message}
}

// MountError represents an error during mount operations
type MountError struct {
	message string
	cause   error
}

func (e *MountError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("mount error: %s (caused by: %v)", e.message, e.cause)
	}
	return fmt.Sprintf("mount error: %s", e.message)
}

// NewMountError creates a new MountError
func NewMountError(message string, cause error) *MountError {
	return &MountError{
		message: message,
		cause:   cause,
	}
} 