// Package utils provides utility functions for the GlusterFS plugin.
// It includes functionality for Unix socket communication and system logging.
package utils

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"glusterfs-plugin/pkg/volume"
)

const (
	// socketDir is the directory where the plugin's Unix socket will be created.
	// This is the standard location for Docker volume plugins.
	socketDir = "/run/docker/plugins"

	// socketName is the name of the Unix socket file.
	// This is the name that Docker will use to communicate with the plugin.
	socketName = "glusterfs.sock"
)

// StartUnixSocket starts the Unix socket server for the volume driver.
// It creates and manages the Unix socket that Docker uses to communicate
// with the plugin.
//
// The function:
// 1. Creates the socket directory if it doesn't exist
// 2. Removes any existing socket file
// 3. Creates a new Unix socket
// 4. Sets appropriate permissions
// 5. Starts accepting connections
//
// Parameters:
// - driver: The volume driver implementation
// - root: The root directory for volume mounts
//
// Returns:
// - error if the server fails to start, nil otherwise
func StartUnixSocket(driver volume.Driver, root string) error {
	// Ensure the socket directory exists
	if err := os.MkdirAll(socketDir, 0755); err != nil {
		return fmt.Errorf("failed to create socket directory: %v", err)
	}

	// Remove existing socket if it exists
	socketPath := filepath.Join(socketDir, socketName)
	if err := os.RemoveAll(socketPath); err != nil {
		return fmt.Errorf("failed to remove existing socket: %v", err)
	}

	// Create Unix socket listener
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return fmt.Errorf("failed to create Unix socket: %v", err)
	}
	defer listener.Close()

	// Set socket permissions
	if err := os.Chmod(socketPath, 0660); err != nil {
		return fmt.Errorf("failed to set socket permissions: %v", err)
	}

	log.Printf("Starting Unix socket server at %s", socketPath)

	// Start accepting connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go handleConnection(conn, driver)
	}
}

// handleConnection handles a single client connection.
// It processes incoming requests from Docker and forwards them to the driver.
//
// The function:
// 1. Reads the request from the connection
// 2. Parses the request
// 3. Calls the appropriate driver method
// 4. Writes the response back to the connection
//
// Parameters:
// - conn: The network connection to handle
// - driver: The volume driver implementation
func handleConnection(conn net.Conn, driver volume.Driver) {
	defer conn.Close()

	// TODO: Implement the actual request handling logic
	// This would involve:
	// 1. Reading the request from the connection
	// 2. Parsing the request
	// 3. Calling the appropriate driver method
	// 4. Writing the response back to the connection
	log.Printf("New connection from %s", conn.RemoteAddr())
} 