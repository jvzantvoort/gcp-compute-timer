package main

import (
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
	"fmt"
)

// SecondsToHuman translate seconds into a human readable format.
func SecondsToHuman(seconds int64) string {
	var retv string

	if seconds > 86400 {
		retv = retv + fmt.Sprintf("%3d days ", seconds/86400)
		seconds = seconds % 86400
	} else {
		retv = retv + "         "
	}

	if seconds > 3600 {
		retv = retv + fmt.Sprintf("%02d:", seconds/3600)
		seconds = seconds % 3600
	} else {
		retv = retv + "00:"
	}

	if seconds > 60 {
		retv = retv + fmt.Sprintf("%02d:", seconds/60)
		seconds = seconds % 60
	} else {
		retv = retv + "00:"
	}

	retv = retv + fmt.Sprintf("%02d", seconds)

	return retv
}

// targetExists return true if target exists
func targetExists(targetpath string) bool {
	_, err := os.Stat(targetpath)
	if err != nil {
		return false
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// Which return path
func Which(command string) string {
	Path := strings.Split(os.Getenv("PATH"), ":")
	var retv string
	for _, dirname := range Path {
		fullpath := path.Join(dirname, command)
		if targetExists(fullpath) {
			retv = fullpath
			break
		}
	}
	return retv
}

// ExitOnError check error and exit if not nil
func ExitOnError(err error) {
	if err != nil {
		log.Errorf("error %v\n", err)
		os.Exit(1)
	}
}
