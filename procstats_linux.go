// +build linux

package sysstats

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ProcStats represents the processes statistics (NOT counted since boot)
type ProcStats struct {
	Running  uint64 // # of processes in runnable state (Linux 2.5.45 onward)
	Blocked  uint64 // # of processes blocked waiting for I/O to complete (Linux 2.5.45 onward)
	RunQueue uint64 // # of currently runnable kernel scheduling entities (processes, threads)
	Total    uint64 // # of kernel scheduling entities that currently exist on the system
}

// ProcRawStats represents the raw processes statistics
type ProcRawStats struct {
	Processes uint64 // # of forks since boot
	ProcStats
	Time int64 // Time when the sample was taken (Unix time)
}

// ProcAvgStats represents the processes statistics
type ProcAvgStats struct {
	NewProcs float64 // # of forks per second
	ProcStats
}

// getProcRawStats gets the processes stats of a linux system from the files
// /proc/loadavg and /proc/stat.
// It returns a ProcRawStats var.
func getProcRawStats() (procRawStats ProcRawStats, err error) {
	procRawStats = ProcRawStats{}

	now := time.Now().Unix()
	procRawStats.Time = now

	// Get runnable and total processes from /proc/loadavg
	loadavg, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return ProcRawStats{}, err
	}
	// Check number of fields in /proc/loadavg
	fields := strings.Fields(strings.TrimSpace(string(loadavg)))
	if len(fields) != 5 {
		return ProcRawStats{}, errors.New("Error parsing file /proc/loadavg. It should have 5 fields")
	}
	// The two values we are interested in are in the fourth field (it consists
	// of two numbers separated by a slash '/')
	field := fields[3]
	fourthField := strings.Split(field, `/`)
	runQueue, err := strconv.ParseUint(fourthField[0], 10, 64)
	procRawStats.RunQueue = runQueue
	total, err := strconv.ParseUint(fourthField[1], 10, 64)
	procRawStats.Total = total

	// Get total, running and blocked processes from /proc/stat
	file, err := os.Open("/proc/stat")
	if err != nil {
		return ProcRawStats{}, err
	}
	defer file.Close()

	reProcs := regexp.MustCompile(`^processes\s+(\d+)`)
	reProcsRunning := regexp.MustCompile(`^procs_running\s+(\d+)`)
	reProcsBlocked := regexp.MustCompile(`^procs_blocked\s+(\d+)`)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if stat := reProcs.FindStringSubmatch(line); stat != nil {
			procs, err := strconv.ParseUint(stat[1], 10, 64)
			if err != nil {
				return ProcRawStats{}, err
			}
			procRawStats.Processes = procs
		} else if stat := reProcsRunning.FindStringSubmatch(line); stat != nil {
			procsRunning, err := strconv.ParseUint(stat[1], 10, 64)
			if err != nil {
				return ProcRawStats{}, err
			}
			procRawStats.Running = procsRunning
		} else if stat := reProcsBlocked.FindStringSubmatch(line); stat != nil {
			procsBlocked, err := strconv.ParseUint(stat[1], 10, 64)
			if err != nil {
				return ProcRawStats{}, err
			}
			procRawStats.Blocked = procsBlocked
		}
	}

	return procRawStats, nil
}

// getProcAvgStats calculates the average between 2 ProcRawStats samples.
func getProcAvgStats(firstSample ProcRawStats, secondSample ProcRawStats) (procAvgStats ProcAvgStats, err error) {
	procAvgStats = ProcAvgStats{}

	timeDelta := float64(secondSample.Time - firstSample.Time)

	// Calculate number of new processes created per second
	if timeDelta > 0 {
		avg := float64(secondSample.Processes-firstSample.Processes) / timeDelta
		procAvgStats.NewProcs = avg
	} else {
		procAvgStats.NewProcs = 0
	}

	// The other values of procAvgStats will be taken from the second sample because
	// they are "current" values (not counted since boot)
	procAvgStats.Running = secondSample.Running
	procAvgStats.Blocked = secondSample.Blocked
	procAvgStats.RunQueue = secondSample.RunQueue
	procAvgStats.Total = secondSample.Total

	return procAvgStats, nil
}

// getProcStatsInterval returns the processes statistics between 2 samples.
// Time interval between the 2 samples is given in seconds
func getProcStatsInterval(interval int64) (procAvgStats ProcAvgStats, err error) {
	firstSample, err := getProcRawStats()
	if err != nil {
		return ProcAvgStats{}, err
	}

	time.Sleep(time.Duration(interval) * time.Second)

	secondSample, err := getProcRawStats()
	if err != nil {
		return ProcAvgStats{}, err
	}

	procAvgStats, err = getProcAvgStats(firstSample, secondSample)
	if err != nil {
		return ProcAvgStats{}, err
	}

	return procAvgStats, nil
}
