package main

import (
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

type Instance struct {
	name               string
	status             string
	LastStartTimestamp string
	starttime          int
}

type Instances struct {
	project   string
	zone      string
	Instances []Instance
}

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
		var inst Instance
		for _, instance := range page.Items {
			log.Println(instance.Name)
			inst.name = instance.Name
			inst.status = instance.Status
			i.Instances = append(i.Instances, inst)
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}

}

func NewInstances(project, zone string) *Instances {
	retv := &Instances{}
	retv.project = project
	retv.zone = zone
	retv.getInstances()
	return retv
}
