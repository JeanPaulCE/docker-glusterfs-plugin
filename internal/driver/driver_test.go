package driver

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"glusterfs-plugin/pkg/types"
	"glusterfs-plugin/pkg/volume"
)

func TestNewDriver(t *testing.T) {
	tests := []struct {
		name    string
		servers []string
		want    *GFSDriver
	}{
		{
			name:    "empty servers",
			servers: []string{},
			want:    &GFSDriver{GFSDriver: types.NewGFSDriver([]string{})},
		},
		{
			name:    "single server",
			servers: []string{"server1"},
			want:    &GFSDriver{GFSDriver: types.NewGFSDriver([]string{"server1"})},
		},
		{
			name:    "multiple servers",
			servers: []string{"server1", "server2"},
			want:    &GFSDriver{GFSDriver: types.NewGFSDriver([]string{"server1", "server2"})},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDriver(tt.servers)
			assert.Equal(t, tt.want.Servers, got.Servers)
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		driver  *GFSDriver
		req     *volume.CreateRequest
		wantErr bool
	}{
		{
			name:    "nil request",
			driver:  NewDriver([]string{}),
			req:     nil,
			wantErr: true,
		},
		{
			name:   "valid servers in env",
			driver: NewDriver([]string{"server1"}),
			req: &volume.CreateRequest{
				Name:    "test",
				Options: map[string]string{},
			},
			wantErr: false,
		},
		{
			name:   "valid servers in options",
			driver: NewDriver([]string{}),
			req: &volume.CreateRequest{
				Name:    "test",
				Options: map[string]string{"servers": "server1"},
			},
			wantErr: false,
		},
		{
			name:   "valid glusteropts",
			driver: NewDriver([]string{}),
			req: &volume.CreateRequest{
				Name:    "test",
				Options: map[string]string{"glusteropts": "-s server1"},
			},
			wantErr: false,
		},
		{
			name:   "invalid: both servers and glusteropts",
			driver: NewDriver([]string{}),
			req: &volume.CreateRequest{
				Name:    "test",
				Options: map[string]string{"servers": "server1", "glusteropts": "-s server1"},
			},
			wantErr: true,
		},
		{
			name:   "invalid: no servers specified",
			driver: NewDriver([]string{}),
			req: &volume.CreateRequest{
				Name:    "test",
				Options: map[string]string{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.driver.Validate(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMountOptions(t *testing.T) {
	tests := []struct {
		name    string
		driver  *GFSDriver
		req     *volume.CreateRequest
		want    []string
	}{
		{
			name:    "nil request",
			driver:  NewDriver([]string{}),
			req:     nil,
			want:    nil,
		},
		{
			name:   "servers from env",
			driver: NewDriver([]string{"server1", "server2"}),
			req: &volume.CreateRequest{
				Name:    "test",
				Options: map[string]string{},
			},
			want: []string{"-s", "server1", "-s", "server2", "--volfile-id=test", "--logger=syslog"},
		},
		{
			name:   "servers from options",
			driver: NewDriver([]string{}),
			req: &volume.CreateRequest{
				Name:    "test",
				Options: map[string]string{"servers": "server1,server2"},
			},
			want: []string{"-s", "server1", "-s", "server2", "--volfile-id=test", "--logger=syslog"},
		},
		{
			name:   "glusteropts",
			driver: NewDriver([]string{}),
			req: &volume.CreateRequest{
				Name:    "test",
				Options: map[string]string{"glusteropts": "-s server1 --volfile-id=test"},
			},
			want: []string{"-s", "server1", "--volfile-id=test", "--logger=syslog"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.driver.MountOptions(tt.req)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPreMount(t *testing.T) {
	// Crear un directorio temporal para las pruebas
	tmpDir, err := os.MkdirTemp("", "glusterfs-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name    string
		driver  *GFSDriver
		req     *volume.MountRequest
		wantErr bool
	}{
		{
			name:    "nil request",
			driver:  NewDriver([]string{}),
			req:     nil,
			wantErr: true,
		},
		{
			name:   "valid mountpoint",
			driver: NewDriver([]string{}),
			req: &volume.MountRequest{
				Name:       "test",
				Mountpoint: filepath.Join(tmpDir, "test"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Crear el directorio de montaje si no existe
			if tt.req != nil {
				if err := os.MkdirAll(tt.req.Mountpoint, 0755); err != nil {
					t.Fatalf("Failed to create mount point: %v", err)
				}
			}

			err := tt.driver.PreMount(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAppendVolumeOptionsByVolumeName(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		volumeName string
		want      []string
	}{
		{
			name:      "empty volume name",
			args:      []string{"mount"},
			volumeName: "",
			want:      []string{"mount"},
		},
		{
			name:      "simple volume",
			args:      []string{"mount"},
			volumeName: "simplevolume",
			want:      []string{"mount", "--volfile-id=simplevolume"},
		},
		{
			name:      "one level subdir",
			args:      []string{"mount"},
			volumeName: "simplevolume/levelone",
			want:      []string{"mount", "--volfile-id=simplevolume", "--subdir-mount=/levelone"},
		},
		{
			name:      "two levels subdir",
			args:      []string{"mount"},
			volumeName: "simplevolume/levelone/level2",
			want:      []string{"mount", "--volfile-id=simplevolume", "--subdir-mount=/levelone/level2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := appendVolumeOptionsByVolumeName(tt.args, tt.volumeName)
			assert.Equal(t, tt.want, got)
		})
	}
} 