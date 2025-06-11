// Package types defines the core data structures used throughout the GlusterFS plugin.
// These types are used to represent the state and configuration of the plugin.
package types

import (
	"glusterfs-plugin/pkg/volume"
)

// Driver interface defines the methods that must be implemented by the GlusterFS driver
type Driver interface {
	Validate(req *volume.CreateRequest) error
	MountOptions(req *volume.CreateRequest) []string
	PreMount(req *volume.MountRequest) error
	PostMount(req *volume.MountRequest)
}

// GFSDriver implements the Driver interface for GlusterFS volumes.
// It provides the core functionality for managing GlusterFS volumes,
// including server configuration and volume operations.
type GFSDriver struct {
	// Servers contains the list of GlusterFS servers to use for volume operations.
	// These servers are used to mount volumes and perform other GlusterFS operations.
	Servers []string
	volume.Driver
}

// NewGFSDriver creates a new instance of the GlusterFS driver.
// It initializes the driver with the provided list of GlusterFS servers.
//
// Parameters:
// - servers: List of GlusterFS server addresses (e.g., ["server1:24007", "server2:24007"])
//
// Returns:
// - A new GFSDriver instance configured with the provided servers
func NewGFSDriver(servers []string) *GFSDriver {
	return &GFSDriver{
		Servers: servers,
	}
} 