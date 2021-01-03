package gcpcomputetimer

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	log "github.com/sirupsen/logrus"
)

// InstanceConfig FIXME putin somthing useful
type InstanceConfig struct {
	Project string
	Zone    string
	Name    string
	MaxAge  int
	Action  string
}

// splitParam split a parameter and return a list of strings
func splitParam(parameter string) []string {
	parameter = strings.TrimSpace(parameter)
	return strings.Split(parameter, "/")
}

// Description return the description of an instance
func (ic InstanceConfig) Description() string {
	return fmt.Sprintf("%s/%s/%s", ic.Project, ic.Zone, ic.Name)
}

// calcMaxUptimeSecs calculate the maximum amount of seconds an instanace may
// be up.
func calcMaxUptimeSecs(calcstring string) (int, error) {

	if number, err := strconv.Atoi(calcstring); err == nil {
		return number, nil
	}

	switch calcstring {
	case "workday":
		return 43200, nil // 12 hrs

	case "day":
		return 86400, nil // 24 hrs

	case "week":
		return 604800, nil

	}

	return 0, fmt.Errorf("Invalid number %s", calcstring)

}

// NewInstanceConfig initialize a new InstanceConfig
func NewInstanceConfig(instr string) (*InstanceConfig, error) {
	retv := &InstanceConfig{}

	// remove comments if any
	if strings.Contains(instr, "#") {
		xcols := strings.Split(instr, "#")
		instr = xcols[0]
	}

	// remove leading and trailing spaces
	instr = strings.TrimSpace(instr)

	// split on cols
	cols := strings.Fields(instr)

	if len(cols) == 0 {
		return retv, errors.New("empty line")
	}

	// split on slash (project/zone/name)
	pcols := splitParam(cols[0])

	if len(pcols) != 3 {
		return retv, fmt.Errorf("pattern mismatch %s", instr)
	}

	retv.Project = pcols[0]
	retv.Zone = pcols[1]
	retv.Name = pcols[2]

	// max age in seconds
	cols = cols[1:]
	if len(cols) == 0 {
		return retv, fmt.Errorf("No age defined")
	}

	if maxage, err := calcMaxUptimeSecs(cols[0]); err == nil {
		retv.MaxAge = maxage
	} else {
		return retv, fmt.Errorf("error %v", err)
	}
	log.Debugf("%s max age %d\n", retv.Description(), retv.MaxAge)

	// actions if any
	cols = cols[1:]
	if len(cols) == 0 {
		retv.Action = "None"
	} else {
		retv.Action = strings.Join(cols, " ")
	}
	return retv, nil
}

// InstanceConfigs FIXME putin somthing useful
type InstanceConfigs struct {
	Bucketname      string
	InstanceConfigs []InstanceConfig
}

func (ics *InstanceConfigs) readConfig() {

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	data, err := ics.readFromBucket(client, "gcp-compute-timer.txt")
	if err != nil {
		log.Fatalf("Cannot read object: %v", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		readline := scanner.Text()
		if instanceconfig_object, err := NewInstanceConfig(readline); err == nil {
			log.Debugf("Config %s/%s/%s", instanceconfig_object.Project, instanceconfig_object.Project, instanceconfig_object.Name)
			ics.InstanceConfigs = append(ics.InstanceConfigs, *instanceconfig_object)
		} else {
			log.Errorf("error %s", err)

		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func (ics InstanceConfigs) readFromBucket(client *storage.Client, object string) ([]byte, error) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	rc, err := client.Bucket(ics.Bucketname).Object(object).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (ics *InstanceConfigs) getDefs(imagename string, defaultval int) (int, string) {
	maxage := defaultval
	action := "None"

	for _, ic := range ics.InstanceConfigs {
		if ic.Name == imagename {
			maxage = ic.MaxAge
			action = ic.Action
			log.Debugf("action: %v", action)
			log.Debugf("maxage: %v", maxage)
		}
	}
	return maxage, action
}

// NewInstanceConfigs initialize the InstanceConfigs object
func NewInstanceConfigs(bucketname string) *InstanceConfigs {
	retv := &InstanceConfigs{}
	retv.Bucketname = bucketname
	retv.readConfig()
	return retv
}
