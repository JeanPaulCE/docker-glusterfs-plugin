// Package volume defines the core interfaces and types for volume management in the GlusterFS plugin.
// It provides the foundation for implementing volume drivers and handling volume operations.
package volume

import (
	"fmt"
)

// Driver defines the interface that volume drivers must implement.
// This interface provides the contract for all volume operations including
// validation, mounting, and capability reporting.
type Driver interface {
	// Validate checks if a volume creation request is valid.
	// It should verify all required parameters and their values.
	Validate(req *CreateRequest) error

	// MountOptions returns the mount options that should be used when mounting a volume.
	// These options are specific to the GlusterFS filesystem and include server
	// information and volume configuration.
	MountOptions(req *CreateRequest) []string

	// PreMount performs any necessary operations before mounting a volume.
	// This includes checking mount point existence and permissions.
	PreMount(req *MountRequest) error

	// PostMount performs any necessary operations after mounting a volume.
	// This includes verifying the mount was successful and logging the result.
	PostMount(req *MountRequest)
}

// CreateRequest represents a request to create a new volume.
// It contains the volume name and any additional options needed for creation.
type CreateRequest struct {
	// Name is the unique identifier for the volume
	Name string

	// Options contains key-value pairs of configuration options.
	// Common options include:
	// - servers: comma-separated list of GlusterFS servers
	// - glusteropts: custom GlusterFS mount options
	Options map[string]string
}

// Validate performs validation checks on the create request.
// It ensures all required fields are present and valid.
//
// Returns:
// - error if validation fails, nil otherwise
func (r *CreateRequest) Validate() error {
	if r == nil {
		return fmt.Errorf("create request cannot be nil")
	}
	if r.Name == "" {
		return fmt.Errorf("volume name cannot be empty")
	}
	if r.Options == nil {
		return fmt.Errorf("options cannot be nil")
	}
	return nil
}

// MountRequest represents a request to mount an existing volume.
// It contains the volume name and the mount point where it should be mounted.
type MountRequest struct {
	// Name is the unique identifier of the volume to mount
	Name string

	// Mountpoint is the absolute path where the volume should be mounted
	Mountpoint string
}

// Validate performs validation checks on the mount request.
// It ensures all required fields are present and valid.
//
// Returns:
// - error if validation fails, nil otherwise
func (r *MountRequest) Validate() error {
	if r == nil {
		return fmt.Errorf("mount request cannot be nil")
	}
	if r.Name == "" {
		return fmt.Errorf("volume name cannot be empty")
	}
	if r.Mountpoint == "" {
		return fmt.Errorf("mount point cannot be empty")
	}
	return nil
}

// Volume represents a volume
type Volume struct {
	Name       string
	Mountpoint string
	Status     map[string]interface{}
}

// Capability represents the capabilities of a driver
type Capability struct {
	Scope string
} 