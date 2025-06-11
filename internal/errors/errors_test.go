package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationError(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    string
	}{
		{
			name:    "simple error",
			message: "test error",
			want:    "validation error: test error",
		},
		{
			name:    "empty message",
			message: "",
			want:    "validation error: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewValidationError(tt.message)
			assert.Equal(t, tt.want, err.Error())
		})
	}
}

func TestMountError(t *testing.T) {
	tests := []struct {
		name    string
		message string
		cause   error
		want    string
	}{
		{
			name:    "error with cause",
			message: "mount failed",
			cause:   NewValidationError("invalid config"),
			want:    "mount error: mount failed (caused by: validation error: invalid config)",
		},
		{
			name:    "error without cause",
			message: "mount failed",
			cause:   nil,
			want:    "mount error: mount failed",
		},
		{
			name:    "empty message",
			message: "",
			cause:   nil,
			want:    "mount error: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewMountError(tt.message, tt.cause)
			assert.Equal(t, tt.want, err.Error())
		})
	}
} 