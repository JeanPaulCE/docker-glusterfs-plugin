package utils

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartSyslog(t *testing.T) {
	// Skip this test if rsyslog is not installed
	if _, err := exec.LookPath("rsyslogd"); err != nil {
		t.Skip("rsyslogd not found, skipping test")
	}

	err := StartSyslog()
	assert.NoError(t, err)
} 