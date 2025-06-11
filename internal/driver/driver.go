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

// GFSDriver implements the Driver interface for GlusterFS volumes
type GFSDriver struct {
	*types.GFSDriver
}

// NewDriver creates a new instance of the GlusterFS driver
func NewDriver(servers []string) *GFSDriver {
	return &GFSDriver{
		GFSDriver: types.NewGFSDriver(servers),
	}
}

// Validate validates the create request
func (p *GFSDriver) Validate(req *volume.CreateRequest) error {
	if req == nil {
		return errors.NewValidationError("create request cannot be nil")
	}

	_, serversDefinedInOpts := req.Options["servers"]
	_, glusteroptsInOpts := req.Options["glusteropts"]

	if len(p.servers) > 0 && (serversDefinedInOpts || glusteroptsInOpts) {
		return errors.NewValidationError("SERVERS is set, options are not allowed")
	}
	if serversDefinedInOpts && glusteroptsInOpts {
		return errors.NewValidationError("servers is set, glusteropts are not allowed")
	}
	if len(p.servers) == 0 && !serversDefinedInOpts && !glusteroptsInOpts {
		return errors.NewValidationError("One of SERVERS, driver_opts.servers or driver_opts.glusteropts must be specified")
	}

	return nil
}

// MountOptions returns the mount options for the volume
func (p *GFSDriver) MountOptions(req *volume.CreateRequest) []string {
	if req == nil {
		log.Printf("warning: MountOptions called with nil request")
		return nil
	}

	servers, serversDefinedInOpts := req.Options["servers"]
	glusteropts, _ := req.Options["glusteropts"]

	var args []string

	if len(p.servers) > 0 {
		for _, server := range p.servers {
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

// PreMount performs pre-mount operations
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

// PostMount performs post-mount operations
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

// appendVolumeOptionsByVolumeName appends the command line arguments into the current argument list given the volume name
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