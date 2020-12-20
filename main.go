package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

const gcp_project string = "jdc-development"
const gcp_zone string = "europe-west4-a"

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
	log.Println("start")
	instances := NewInstances(gcp_project, gcp_zone)
	for _, instance := range instances.Instances {
		fmt.Printf("%s %s\n", instance.name, instance.status)
	}

}

// vim: noexpandtab filetype=go
