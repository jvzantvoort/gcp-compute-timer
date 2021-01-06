package gcpcomputetimer

import (
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

func translateAction(instance Instance, noexec bool) {

	log.Debugf("Action[%s]: start", instance.Name)
	defer log.Debugf("Action[%s]: end", instance.Name)

	if instance.Action == "None" {
		log.Debugf("Action[%s]: action is None", instance.Name)
		return
	}

	if instance.Action == "stop" {
		log.Debugf("Action[%s]: action is stop", instance.Name)
		if noexec {
			log.Infof("Action[%s]: Would have stopped instance", instance.Name)
			return
		}
		instance.Stop()
		return
	}


	if strings.HasPrefix(instance.Action, "shell ") {
		action :=  strings.TrimSpace(strings.TrimLeft(instance.Action, "shell"))
		log.Debugf("Action[%s]: shell action: %s", instance.Name, action)
	}

}

func RunCheck(gcp_project string, gcp_zone string, gcp_bucket string, noexec bool) {
	instanceconfig := NewInstanceConfigs(gcp_bucket)
	log.Debugf("RunCheck[%s/%s]: start", gcp_project, gcp_zone)
	defer log.Debugf("RunCheck[%s/%s]: end", gcp_project, gcp_zone)

	// Get the instances their state and if too old
	instances := NewInstances(gcp_project, gcp_zone)
	for indx, instance := range instances.Instances {
		checkInstance(&instance, *instanceconfig)
		instances.Instances[indx] = instance
	}

	if noexec {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Project", "Zone", "Name", "Status", "Too Old", "Uptime", "Action"})
		for _, instance := range instances.Instances {
			var cols []string
			cols = append(cols, instance.Project)
			cols = append(cols, instance.Zone)
			cols = append(cols, instance.Name)
			cols = append(cols, instance.Status)
			if instance.TooOld {
				cols = append(cols, "true")
			} else {
				cols = append(cols, "false")

			}
			cols = append(cols, SecondsToHuman(int64(instance.Age)))
			cols = append(cols, instance.Action)
			table.Append(cols)
		}
		table.SetHeaderLine(true)
		table.SetBorder(false)
		table.Render()
		return
	}

	for _, instance := range instances.Instances {
		log.Debugf("RunCheck[%s/%s/%s]: start", gcp_project, gcp_zone, instance.Name)
		if instance.TooOld {
			log.Warnf("RunCheck[%s/%s/%s]: instance is too old", gcp_project, gcp_zone, instance.Name)
			translateAction(instance, noexec)
		} else {
			log.Debugf("RunCheck[%s/%s/%s]: instance is not too old", gcp_project, gcp_zone, instance.Name)

		}
		log.Debugf("RunCheck[%s/%s/%s]: end", gcp_project, gcp_zone, instance.Name)

	}
}
