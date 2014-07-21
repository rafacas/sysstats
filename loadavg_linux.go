// +build linux

package sysstats

import (
	"io/ioutil"
	"strconv"
	"strings"
)

// LoadAvg represents the load average of the system
type LoadAvg struct {
	Avg1  float64 // The average processor workload of the last minute
	Avg5  float64 // The average processor workload of the last 5 minutes
	Avg15 float64 // The average processor workload of the last 15 minutes
}

// getLoadAvg gets the load average of a linux system from the
// file /proc/loadavg.
func getLoadAvg() (loadAvg *LoadAvg, err error) {
	file, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return nil, err
	}
	content := string(file[:len(file)])

	loadAvg = new(LoadAvg)
	fields := strings.Fields(content)
	for i := 0; i < 3; i++ {
		load, err := strconv.ParseFloat(fields[i], 64)
		if err != nil {
			return nil, err
		}
		switch i {
		case 0:
			loadAvg.Avg1 = load
		case 1:
			loadAvg.Avg5 = load
		case 2:
			loadAvg.Avg15 = load
		}
	}

	return loadAvg, nil
}
