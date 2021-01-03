package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	gcpt "github.com/jvzantvoort/gcp-compute-timer"
	"github.com/jvzantvoort/gcp-compute-timer/config"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

}

func main() {

	configuration := config.NewConfiguration()

	gcp_project := configuration.GCP.Project
	gcp_zone := configuration.GCP.Zone
	gcp_bucket := configuration.GCP.Bucket

	gcpt.Worker(gcp_project, gcp_zone, gcp_bucket)
}

// vim: noexpandtab filetype=go
