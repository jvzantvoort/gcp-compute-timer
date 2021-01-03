package gcpcomputetimer

import (
	log "github.com/sirupsen/logrus"
)

const (
	DEFAULT_MAX_AGE int = 43200
)

func Worker(project_name, zone_name, bucket_name string) {
	inst_cfg := NewInstanceConfigs(bucket_name)

	// Get the instances and their state
	instances := NewInstances(project_name, zone_name)
	for _, instance := range instances.Instances {
		var action string
		instance.maxage, action = inst_cfg.getDefs(instance.name, DEFAULT_MAX_AGE)

		if instance.status == "RUNNING" {
			if instance.IsTooOld() {
				log.Warningf("image: %s state: %s age : %s\n", instance.name, instance.status, SecondsToHuman(instance.age))
				log.Warningf("  max age: %s\n", SecondsToHuman(int64(instance.maxage)))
				log.Warningf("  action: %s", action)

			} else {
				log.Infof("image: %s state: %s age : %s\n", instance.name, instance.status, SecondsToHuman(instance.age))
			}

		} else {
			log.Infof("image: %s state: %s\n", instance.name, instance.status)
		}
	}

}

// vim: noexpandtab filetype=go
