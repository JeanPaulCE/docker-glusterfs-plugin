package types

import (
	"github.com/docker/go-plugins-helpers/volume"
	mountedvolume "github.com/marcelo-ochoa/docker-volume-plugins/mounted-volume"
)

// Driver interface defines the methods that must be implemented by the GlusterFS driver
type Driver interface {
	Validate(req *volume.CreateRequest) error
	MountOptions(req *volume.CreateRequest) []string
	PreMount(req *volume.MountRequest) error
	PostMount(req *volume.MountRequest)
}

// GFSDriver represents a GlusterFS volume driver
type GFSDriver struct {
	servers []string
	mountedvolume.Driver
}

// NewGFSDriver creates a new instance of the GlusterFS driver
func NewGFSDriver(servers []string) *GFSDriver {
	return &GFSDriver{
		Driver:  *mountedvolume.NewDriver("glusterfs", true, "gfs", "local"),
		servers: servers,
	}
} 