package config

import (
	"fmt"
	"log"
	"os/user"
	"path"

	"github.com/spf13/viper"
)

// GCPConfig structure containing:
// Project name of the project
// Zone    zone the instances are located
// Bucket  bucket used in accounting
type GCPConfig struct {
	Project string
	Zone    string
	Bucket  string
}

const (
	constMainConfig string = "/etc/gcp"
	constMainName   string = "gcp-compute-timer"
)

// Configuration container for configuration. Currently only contains GCP.
type Configuration struct {
	GCP GCPConfig
}

// NewConfiguration initialize a new Configuration object.
func NewConfiguration() Configuration {
	var retv Configuration

	// Get user info
	usrobj, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	configdir := path.Join(usrobj.HomeDir, ".config")
	viper.SetConfigName(constMainName)   // name of config file (without extension)
	viper.AddConfigPath(configdir)       // optionally look for config in the ~/.config directory
	viper.AddConfigPath(constMainConfig) // optionally look for config in the working directory

	err = viper.ReadInConfig() // Find and read the config file

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	err = viper.Unmarshal(&retv)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	return retv
}
