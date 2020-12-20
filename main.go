package main

import (
	"fmt"
	"os"
	"os/user"
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/jvzantvoort/gcp-compute-timer/config"
)

const (
	DEFAULT_MAX_AGE int = 43200
)

func getconfiguration(homedir string) config.Configuration {
	var configuration config.Configuration
	configdir := path.Join(homedir, ".config")
	viper.SetConfigName("gcp-compute-timer") // name of config file (without extension)
	viper.AddConfigPath(configdir)           // optionally look for config in the ~/.config directory
	viper.AddConfigPath("/etc/gcp")          // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	err = viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	return configuration
}

func main() {

	// Get user info
	usrobj, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	homedir := usrobj.HomeDir
	configuration := getconfiguration(homedir)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	gcp_project := configuration.GCP.Project
	gcp_zone := configuration.GCP.Zone

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

	nisconf := NewInstanceConfigs("data.txt")

	// Get the instances and their state
	instances := NewInstances(gcp_project, gcp_zone)
	for _, instance := range instances.Instances {
		var action string
		instance.maxage, action = nisconf.getDefs(instance.name, DEFAULT_MAX_AGE)

		if instance.status == "RUNNING" {
			if instance.IsTooOld() {
				log.Warningf("image: %s state: %s age : %s\n", instance.name, instance.status, SecondsToHuman(instance.age))
				log.Debugf("action: %s", action)
			} else {
				log.Infof("image: %s state: %s age : %s\n", instance.name, instance.status, SecondsToHuman(instance.age))
			}

		} else {
			log.Infof("image: %s state: %s\n", instance.name, instance.status)
		}
	}
}

// vim: noexpandtab filetype=go
