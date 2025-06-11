package utils

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartSyslog(t *testing.T) {
	// Skip this test if rsyslog is not installed
	if _, err := exec.LookPath("rsyslogd"); err != nil {
		t.Skip("rsyslogd not found, skipping test")
	}

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "start syslog",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StartSyslog()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
} 