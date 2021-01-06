// Monitor the uptime of gcp computer instances and act accordingly.
//
package gcpcomputetimer

import (
	log "github.com/sirupsen/logrus"
)

// checkInstance check an instance and return
//
func checkInstance(instance *Instance, instanceconfig InstanceConfigs) {

	log.Debugf("checkInstance[%s]: start", instance.Name)
	defer log.Debugf("checkInstance[%s]: end", instance.Name)

	instance.MaxAge, instance.Action = instanceconfig.getDefs(instance.Name, constMaxAge)
	// Translate shell commands
	instance.Action = instance.Parse(instance.Action)

	// Check if it is too old
	if instance.Status != "RUNNING" {
		instance.TooOld = false
		instance.Age = 0
	} else {
		if int(instance.Age) > instance.MaxAge {
			instance.TooOld = true
		} else {
			instance.TooOld = false
		}
	}
}

// vim: noexpandtab filetype=go
