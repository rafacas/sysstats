// +build linux

package sysstats

import (
	"bufio"
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
//  User      - % Cpu time spent in user mode.
//  Nice      - % Cpu time spent in user mode with low priority (nice).
//  System    - % Cpu time spent in system mode.
//  Idle      - % Cpu time spent in the idle task.
//  Iowait    - % Cpu time spent waiting for I/O to complete (since 2.5.41).
//  Irq       - % Cpu time servicing interrupts (since 2.6.0-test4).
//  Softirq   - % Cpu time servicing softirqs (since 2.6.0-test4).
//  Steal     - % Cpu stolen time, which is the time spent in other operating
//              systems when running a virtualized environment (since 2.6.11).
//  Guest     - % Cpu time spent running a virtual Cpu for guest operating
//              systems under the control of the Linux kernel (since 2.6.24).
//  GuestNice - % Cpu time spent running a niced guest (virtual Cpu for guest
//              operating systems under the control of the Linux kernel)
//              (since 2.6.33).
//  Total     - % Cpu utilization (not idle).
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

func getCpuStats() (cpuStats CpusRawStats, err error) {
	return getCpuRawStats()
}
