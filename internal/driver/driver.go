// Package driver implements the GlusterFS volume driver functionality.
// It provides the concrete implementation of the volume.Driver interface
// for managing GlusterFS volumes in Docker.
package driver

import (
	"fmt"
	"log"
	"os"
	"strings"
	"glusterfs-plugin/internal/errors"
	"glusterfs-plugin/pkg/types"
	"glusterfs-plugin/pkg/volume"
)

// GFSDriver implements the Driver interface for GlusterFS volumes.
// It wraps the base GFSDriver from pkg/types and adds Docker-specific functionality.
type GFSDriver struct {
	*types.GFSDriver
}

// NewDriver creates a new instance of the GlusterFS driver.
// It initializes the driver with the provided list of GlusterFS servers.
//
// Parameters:
// - servers: List of GlusterFS server addresses
//
// Returns:
// - A new GFSDriver instance configured with the provided servers
func NewDriver(servers []string) *GFSDriver {
	return &GFSDriver{
		GFSDriver: types.NewGFSDriver(servers),
	}
}

// Validate validates the create request.
// It ensures that the request is valid and that the server configuration
// is consistent with the provided options.
//
// The validation rules are:
// 1. If SERVERS is set, no options are allowed
// 2. If servers is set in options, glusteropts are not allowed
// 3. At least one of SERVERS, servers, or glusteropts must be specified
//
// Parameters:
// - req: The create request to validate
//
// Returns:
// - error if validation fails, nil otherwise
func (p *GFSDriver) Validate(req *volume.CreateRequest) error {
	if req == nil {
		return errors.NewValidationError("create request cannot be nil")
	}

	_, serversDefinedInOpts := req.Options["servers"]
	_, glusteroptsInOpts := req.Options["glusteropts"]

	if len(p.Servers) > 0 && (serversDefinedInOpts || glusteroptsInOpts) {
		return errors.NewValidationError("SERVERS is set, options are not allowed")
	}
	if serversDefinedInOpts && glusteroptsInOpts {
		return errors.NewValidationError("servers is set, glusteropts are not allowed")
	}
	if len(p.Servers) == 0 && !serversDefinedInOpts && !glusteroptsInOpts {
		return errors.NewValidationError("One of SERVERS, driver_opts.servers or driver_opts.glusteropts must be specified")
	}

	return nil
}

// MountOptions returns the mount options for the volume.
// It constructs the appropriate mount options based on the server configuration
// and volume options.
//
// The mount options include:
// - Server addresses (-s option)
// - Volume ID (--volfile-id)
// - Subdirectory mount point (--subdir-mount) if specified
// - Logger configuration (--logger=syslog)
//
// Parameters:
// - req: The create request containing volume options
//
// Returns:
// - List of mount options to use when mounting the volume
func (p *GFSDriver) MountOptions(req *volume.CreateRequest) []string {
	if req == nil {
		log.Printf("warning: MountOptions called with nil request")
		return nil
	}

	servers, serversDefinedInOpts := req.Options["servers"]
	glusteropts, _ := req.Options["glusteropts"]

	var args []string

	if len(p.Servers) > 0 {
		for _, server := range p.Servers {
			args = append(args, "-s", server)
		}
		args = appendVolumeOptionsByVolumeName(args, req.Name)
	} else if serversDefinedInOpts {
		for _, server := range strings.Split(servers, ",") {
			args = append(args, "-s", server)
		}
		args = appendVolumeOptionsByVolumeName(args, req.Name)
	} else {
		args = strings.Split(glusteropts, " ")
	}

	args = append(args, "--logger=syslog")
	return args
}

// PreMount performs pre-mount operations.
// It verifies that the mount point exists and is accessible.
//
// Parameters:
// - req: The mount request containing the mount point
//
// Returns:
// - error if pre-mount checks fail, nil otherwise
func (p *GFSDriver) PreMount(req *volume.MountRequest) error {
	if req == nil {
		return errors.NewMountError("mount request cannot be nil", nil)
	}

	// Verify that the mount point exists and is accessible
	if _, err := os.Stat(req.Mountpoint); err != nil {
		return errors.NewMountError(
			fmt.Sprintf("mount point %s is not accessible", req.Mountpoint),
			err,
		)
	}

	return nil
}

// PostMount performs post-mount operations.
// It verifies that the mount was successful and logs the result.
//
// Parameters:
// - req: The mount request containing the mount point
func (p *GFSDriver) PostMount(req *volume.MountRequest) {
	if req == nil {
		log.Printf("warning: PostMount called with nil request")
		return
	}

	// Verify that the mount was successful
	if _, err := os.Stat(req.Mountpoint); err != nil {
		log.Printf("error: Mount point %s is not accessible after mount: %v", req.Mountpoint, err)
		return
	}

	log.Printf("successfully mounted volume %s at %s", req.Name, req.Mountpoint)
}

// appendVolumeOptionsByVolumeName appends the command line arguments for volume mounting.
// It handles both simple volume names and subdirectory mounts.
//
// The function adds:
// - --volfile-id for the volume name
// - --subdir-mount for any subdirectory path
//
// Parameters:
// - args: The existing command line arguments
// - volumeName: The name of the volume, optionally including a subdirectory path
//
// Returns:
// - Updated list of command line arguments
func appendVolumeOptionsByVolumeName(args []string, volumeName string) []string {
	if volumeName == "" {
		log.Printf("warning: appendVolumeOptionsByVolumeName called with empty volume name")
		return args
	}

	parts := strings.SplitN(volumeName, "/", 2)
	ret := append(args, "--volfile-id="+parts[0])
	if len(parts) == 2 {
		ret = append(ret, "--subdir-mount=/"+parts[1])
	}
	return ret
} 