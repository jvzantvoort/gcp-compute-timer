package gcpcomputetimer

import (
	"time"

	log "github.com/sirupsen/logrus"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

// Instance struct to contain the relevant element of an instance
type Instance struct {
	Name               string
	Project            string
	Zone               string
	Status             string
	LastStartTimestamp string
	TooOld             bool
	Age                int64
	Action             string
	MaxAge             int
}

// Instances struct to contain all the instances
type Instances struct {
	Project   string
	Zone      string
	Now       int64
	Instances []Instance
}

// StartTime converts the LastStartTimestamp string to an epoch int64
func (in *Instance) StartTime() int64 {
	retv, err := time.Parse(time.RFC3339, in.LastStartTimestamp)
	if err != nil {
		log.Error("cannot parse LastStartTimestamp")
	}
	return retv.Unix()
}

func (in Instance) Stop() bool {
	log.Debugf("Stop[%s]: start", in.Name)
	defer log.Debugf("Stop[%s]: stop", in.Name)

	ctx := context.Background()

	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	computeService, err := compute.New(c)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := computeService.Instances.Stop(in.Project, in.Zone, in.Name).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Change code below to process the `resp` object:
	log.Debugf("Stop[%s]: response: %#v", in.Name, resp)
	return true
}

// loadInstances loads the instance information from google
func (i *Instances) loadInstances() {
	ctx := context.Background()

	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	computeService, err := compute.New(c)
	if err != nil {
		log.Fatal(err)
	}

	req := computeService.Instances.List(i.Project, i.Zone)

	if err := req.Pages(ctx, func(page *compute.InstanceList) error {
		for _, instance := range page.Items {
			var inst Instance
			inst.Name = instance.Name
			inst.Status = instance.Status
			inst.Project = i.Project
			inst.Zone = i.Zone
			inst.LastStartTimestamp = instance.LastStartTimestamp
			inst.Age = i.Now - inst.StartTime()
			inst.MaxAge = 86400
			i.Instances = append(i.Instances, inst)
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}

}

// NewInstances initializes the instances object
func NewInstances(project, zone string) *Instances {
	retv := &Instances{}
	now := time.Now()
	retv.Now = now.Unix()
	retv.Project = project
	retv.Zone = zone
	retv.loadInstances()
	return retv
}
