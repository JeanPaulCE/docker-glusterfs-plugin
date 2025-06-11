package main

import (
	"flag"
	"log"
	"strings"

	"glusterfs-plugin/internal/driver"
	"glusterfs-plugin/internal/utils"
)

var (
	servers = flag.String("servers", "", "Comma separated list of GlusterFS servers")
	root    = flag.String("root", "/mnt/glusterfs", "Mount root of volume plugin")
)

func main() {
	flag.Parse()

	if *servers == "" {
		log.Fatal("servers parameter is required")
	}

	serversList := strings.Split(*servers, ",")
	for i, server := range serversList {
		serversList[i] = strings.TrimSpace(server)
	}

	d := driver.NewDriver(serversList)
	if err := utils.StartUnixSocket(d, *root); err != nil {
		log.Fatal(err)
	}
} 