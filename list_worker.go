package gcpcomputetimer

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

func PrintCheck(gcp_project string, gcp_zone string, gcp_bucket string, printfull bool) {
	instanceconfig := NewInstanceConfigs(gcp_bucket)
	log.Debugf("PrintCheck[%s/%s/%s]: start", gcp_project, gcp_zone, gcp_bucket)
	defer log.Debugf("PrintCheck[%s/%s/%s]: end", gcp_project, gcp_zone, gcp_bucket)

	// Get the instances their state and if too old
	instances := NewInstances(gcp_project, gcp_zone)
	for _, instance := range instances.Instances {
		checkInstance(&instance, *instanceconfig)
	}

	if printfull {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Project", "Zone", "Name", "Status", "Too Old"})
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
			table.Append(cols)
		}
		table.SetHeaderLine(true)
		table.SetBorder(false)
		table.Render()
	} else {
		for _, instance := range instances.Instances {
			fmt.Printf("%s\n", instance.Name)
		}
	}
}
