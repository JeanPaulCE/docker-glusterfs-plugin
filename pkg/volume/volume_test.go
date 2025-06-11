package volume

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *CreateRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: &CreateRequest{
				Name:    "test",
				Options: map[string]string{},
			},
			wantErr: false,
		},
		{
			name: "empty name",
			req: &CreateRequest{
				Name:    "",
				Options: map[string]string{},
			},
			wantErr: true,
		},
		{
			name: "nil options",
			req: &CreateRequest{
				Name:    "test",
				Options: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMountRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *MountRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: &MountRequest{
				Name:       "test",
				Mountpoint: "/mnt/test",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			req: &MountRequest{
				Name:       "",
				Mountpoint: "/mnt/test",
			},
			wantErr: true,
		},
		{
			name: "empty mountpoint",
			req: &MountRequest{
				Name:       "test",
				Mountpoint: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
} 