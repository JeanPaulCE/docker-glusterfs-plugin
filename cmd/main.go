package main

import (
	"log"
	"os"
	"strings"

	"github.com/marcelo-ochoa/docker-volume-plugins/glusterfs-plugin/internal/driver"
	"github.com/marcelo-ochoa/docker-volume-plugins/glusterfs-plugin/internal/utils"
)

func main() {
	// Initialize logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("starting glusterfs volume plugin")

	// Start syslog
	if err := utils.StartSyslog(); err != nil {
		log.Fatalf("failed to initialize syslog: %v", err)
	}

	// Build and initialize driver
	var servers []string
	if os.Getenv("SERVERS") != "" {
		servers = strings.Split(os.Getenv("SERVERS"), ",")
	}

	d := driver.NewDriver(servers)

	if os.Getenv("SECURE_MANAGEMENT") != "" {
		file, err := os.Create("/var/lib/glusterd/secure-access")
		if err != nil {
			log.Fatalf("failed to create secure-access file: %v", err)
		}
		defer file.Close()
	}

	d.Init(d)
	defer d.Close()

	// Start serving
	log.Printf("glusterfs volume plugin is ready to serve")
	if err := d.ServeUnix(); err != nil {
		log.Fatalf("failed to serve unix socket: %v", err)
	}
} 