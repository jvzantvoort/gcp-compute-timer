package main

import (
	"time"

	log "github.com/sirupsen/logrus"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

// Struct to contain the relevant element of an instance
type Instance struct {
	name               string
	status             string
	LastStartTimestamp string
	age                int64
	maxage             int
}

// Struct to contain all the instances
type Instances struct {
	project   string
	zone      string
	now       int64
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

func (in Instance) IsTooOld() bool {
	if int(in.age) > in.maxage {
		return true
	}
	return false
}

// getInstances loads the instance information from google
func (i *Instances) getInstances() {
	ctx := context.Background()

	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	computeService, err := compute.New(c)
	if err != nil {
		log.Fatal(err)
	}

	req := computeService.Instances.List(i.project, i.zone)

	if err := req.Pages(ctx, func(page *compute.InstanceList) error {
		for _, instance := range page.Items {
			var inst Instance
			inst.name = instance.Name
			inst.status = instance.Status
			inst.LastStartTimestamp = instance.LastStartTimestamp
			inst.age = i.now - inst.StartTime()
			inst.maxage = 86400
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
	retv.now = now.Unix()
	retv.project = project
	retv.zone = zone
	retv.getInstances()
	return retv
}
