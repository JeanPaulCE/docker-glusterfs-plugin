package utils

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

// StartSyslog starts the rsyslog daemon
func StartSyslog() error {
	cmd := exec.Command("rsyslogd", "-n")
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start rsyslog: %v", err)
	}

	// Wait a short time to ensure rsyslog started properly
	time.Sleep(100 * time.Millisecond)
	
	// Check if rsyslog is running
	if cmd.Process == nil {
		return fmt.Errorf("rsyslog process failed to start")
	}

	log.Printf("rsyslog daemon started successfully")
	return nil
} 