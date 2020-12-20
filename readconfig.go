package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"errors"
	"bufio"
	log "github.com/sirupsen/logrus"
)

type InstanceConfig struct {
	Project string
	Zone    string
	Name    string
	MaxAge  int
	Action  string
}

type InstanceConfigs struct {
	InstanceConfigs []InstanceConfig
}

func splitParam(parameter string) []string {
	parameter = strings.TrimSpace(parameter)
	return strings.Split(parameter, "/")
}

func (ic InstanceConfig) Description() string {
	return fmt.Sprintf("%s/%s/%s", ic.Project, ic.Zone, ic.Name)
}

func calcAge(calcstring string) (int, error) {

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

	if maxage, err := calcAge(cols[0]); err == nil {
		retv.MaxAge = maxage
	} else {
		return retv, fmt.Errorf("error %v\n", err)
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

func (ics *InstanceConfigs) readconfig(filepath string) {

	if targetExists(filepath) {
		file, err := os.Open(filepath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			read_line := scanner.Text()
			if xxxxx, err := NewInstanceConfig(read_line); err == nil {
				log.Debugf("Config %s/%s/%s", xxxxx.Project, xxxxx.Project, xxxxx.Name)
				ics.InstanceConfigs = append(ics.InstanceConfigs, *xxxxx)
			} else {
				log.Errorf("error %s", err)

			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Errorf("Cannot open %s\n", filepath)
	}
}

func (ics *InstanceConfigs) getDefs(imagename string, defaultval int)  (int, string) {
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

func NewInstanceConfigs(filepath string) *InstanceConfigs {
	retv := &InstanceConfigs{}
	retv.readconfig(filepath)
	return retv
}
