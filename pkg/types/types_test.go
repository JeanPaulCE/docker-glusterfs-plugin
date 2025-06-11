package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGFSDriver(t *testing.T) {
	tests := []struct {
		name    string
		servers []string
		want    *GFSDriver
	}{
		{
			name:    "empty servers",
			servers: []string{},
			want:    &GFSDriver{Servers: []string{}},
		},
		{
			name:    "single server",
			servers: []string{"server1"},
			want:    &GFSDriver{Servers: []string{"server1"}},
		},
		{
			name:    "multiple servers",
			servers: []string{"server1", "server2"},
			want:    &GFSDriver{Servers: []string{"server1", "server2"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGFSDriver(tt.servers)
			assert.Equal(t, tt.want.Servers, got.Servers)
		})
	}
} 