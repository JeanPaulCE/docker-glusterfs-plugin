package volume

// CreateRequest represents a request to create a volume
type CreateRequest struct {
	Name    string
	Options map[string]string
}

// MountRequest represents a request to mount a volume
type MountRequest struct {
	Name       string
	Mountpoint string
}

// Driver interface defines the methods that must be implemented by volume drivers
type Driver interface {
	// Create creates a new volume
	Create(req *CreateRequest) error
	// Remove removes a volume
	Remove(req *CreateRequest) error
	// Mount mounts a volume
	Mount(req *MountRequest) error
	// Unmount unmounts a volume
	Unmount(req *MountRequest) error
	// Get gets information about a volume
	Get(req *CreateRequest) (*Volume, error)
	// List lists all volumes
	List() ([]*Volume, error)
	// Path returns the path to a volume
	Path(req *CreateRequest) (string, error)
	// Capabilities returns the capabilities of the driver
	Capabilities() *Capability
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