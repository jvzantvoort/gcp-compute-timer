package gcpcomputetimer

import (
	log "github.com/sirupsen/logrus"
)


func InstanceWorker(instance Instance, instanceconfig InstanceConfigs) {

	instance_name := instance.name
	log.Debugf("%s: start", instance_name)
	defer log.Debugf("%s: end", instance_name)

	var action string
	instance.maxage, action = instanceconfig.getDefs(instance_name, DefaultMaxAge)

	if instance.status != "RUNNING" {
		log.Debugf("%s: state: %s", instance_name, instance.status)
		return
	}

	if !instance.IsTooOld() {
		log.Debugf("%s: state: %s age : %s (not too old)", instance_name, instance.status, SecondsToHuman(instance.age))
		return
	}

	log.Warningf("image: %s state: %s age : %s\n", instance_name, instance.status, SecondsToHuman(instance.age))
	log.Warningf("  max age: %s\n", SecondsToHuman(int64(instance.maxage)))
	log.Warningf("  action: %s", action)

	switch action {
	case "None":
		return
	case "stop":
		instance.Stop()

	default:
		log.Warningf("  action: %s", action)

	}
}

func Worker(project_name, zone_name, bucket_name string) {
	instanceconfig := NewInstanceConfigs(bucket_name)

	// Get the instances and their state
	instances := NewInstances(project_name, zone_name)
	for _, instance := range instances.Instances {
		InstanceWorker(instance, *instanceconfig)
	}

}

// vim: noexpandtab filetype=go
