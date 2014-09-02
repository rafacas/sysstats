// +build linux

package sysstats

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// CpuRawStats represents *one* CPU raw statistics of a linux system.
type CpuRawStats map[string]uint64

// CpuStats represents *one* CPU statistics of a linux system.
type CpuStats map[string]float64

// CpusRawStats represents *all* the CPU raw statistics of a linux system.
type CpusRawStats map[string]CpuRawStats

// CpusStats represents *all* the CPU statistics of a linux system.
type CpusStats map[string]CpuStats

// getCpuRawStats gets the Cpu raw stats of a linux system from the
// file /proc/stat
// It retuns a CpusRawStats var. CpusRawStats is a map which keys are the
// 'kernel name' of the CPUs (cpu, cpu0, etc). The value of each key is a
// CpuRawStats var with the statistics for that CPU.
// CpusRawStats has the following key:
//  Name      - Name of the Cpu (as it is on /proc/stat: cpu, cpu0,...).
// CpuRawStats has the following keys:
//  User      - Time spent in user mode.
//  Nice      - Time spent in user mode with low priority (nice).
//  System    - Time spent in system mode.
//  Idle      - Time spent in the idle task.
//  Iowait    - Time spent waiting for I/O to complete (since 2.5.41).
//  Irq       - Time servicing interrupts (since 2.6.0-test4).
//  Softirq   - Time servicing softirqs (since 2.6.0-test4).
//  Steal     - Stolen time, which is the time spent in other operating
//              systems when running a virtualized environment (since 2.6.11).
//  Guest     - Time spent running a virtual Cpu for guest operating
//              systems under the control of the Linux kernel (since 2.6.24).
//  GuestNice - Time spent running a niced guest (virtual Cpu for guest
//              operating systems under the control of the Linux kernel)
//              (since 2.6.33).
//  Total     - Total time.
func getCpuRawStats() (cpusRawStats CpusRawStats, err error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cpusRawStats = CpusRawStats{}

	re := regexp.MustCompile(`^cpu.*$`)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		stats := re.FindString(line)
		if stats == "" {
			// No match so no more cpu 'lines'
			break
		}
		fmt.Println(stats)
		cpuName, rawStats, err := parseCpuRawStats(stats)
		if err != nil {
			return nil, err
		}
		cpusRawStats[cpuName] = rawStats
	}

	return cpusRawStats, nil
}

// parseCpuRawStats parses the CPU stats as they are in the file /proc/stat
// stats has the following format:
//   cpu  294 0 309 10612 71 30 0 0 0 0
// cpuName is the name of the CPU (cpu, cpu0, cpu1, etc) and rawStats has
// the following format:
//   map[User:9366 Nice:0 System:5692 Iowait:114 Steal:0 GuestNice:0
//       Idle:1458880 Irq:806 Softirq:0 Guest:0]
func parseCpuRawStats(stats string) (cpuName string, rawStats CpuRawStats,
	err error) {
	rawStats = CpuRawStats{}

	fields := strings.Fields(stats)
	cpuName = fields[0]
	for i := 1; i < len(fields); i++ {
		stat, err := strconv.ParseUint(fields[i], 10, 64)
		if err != nil {
			return "", nil, err
		}
		rawStats[`Total`] += stat
		switch i {
		case 1:
			rawStats[`User`] = stat
		case 2:
			rawStats[`Nice`] = stat
		case 3:
			rawStats[`System`] = stat
		case 4:
			rawStats[`Idle`] = stat
		case 5:
			rawStats[`Iowait`] = stat
		case 6:
			rawStats[`Irq`] = stat
		case 7:
			rawStats[`Softirq`] = stat
		case 8:
			rawStats[`Steal`] = stat
		case 9:
			rawStats[`Guest`] = stat
		case 10:
			rawStats[`GuestNice`] = stat
		}
	}

	return cpuName, rawStats, nil
}

// getCpuStats calculates the average between 2 CpusRawStats samples and returns
// the % CPU utilization
func getCpuStats(firstSample CpusRawStats, secondSample CpusRawStats) (cpusStats CpusStats, err error) {
	fmt.Println(firstSample)
	fmt.Println(secondSample)

	cpusStats = CpusStats{}

	for cpuName, secondRawStats := range secondSample {
		matched, err := regexp.MatchString(`^cpu.*$`, cpuName)
		if err != nil {
			return nil, err
		}
		if !matched {
			return nil, errors.New("cpuName doesn't match the pattern")
		}

		firstRawStats, ok := firstSample[cpuName]
		if !ok {
			return nil, errors.New("The key " + cpuName + " doesn't exist in the first sample of CpusRawStats")
		}

		cpuStats := CpuStats{}
		totalTime := float64(secondRawStats[`Total`] - firstRawStats[`Total`])
		// Calculate the average between the two samples
		for key, secondValue := range secondRawStats {
			// Don't calculate the average if the key is 'Total'
			if key == `Total` {
				continue
			}
			avg := float64(secondValue-firstRawStats[key]) * 100.00 / totalTime
			avgStr := fmt.Sprintf("%3.2f", avg)
			cpuStats[key], err = strconv.ParseFloat(avgStr, 64)
			if err != nil {
				return nil, err
			}

		}
		cpuTotal := 100.00 - cpuStats[`Idle`]
		cpuTotalStr := fmt.Sprintf("%3.2f", cpuTotal)
		cpuStats[`Total`], err = strconv.ParseFloat(cpuTotalStr, 64)
		if err != nil {
			return nil, err
		}

		cpusStats[cpuName] = cpuStats
	}

	return cpusStats, nil
}
