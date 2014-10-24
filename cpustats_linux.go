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
	"time"
)

// CpuRawStats represents *one* CPU raw statistics of a linux system.
// The amount of time, measured in units of USER_HZ
// Keys:
//   User      - Time spent in user mode.
//   Nice      - Time spent in user mode with low priority (nice).
//   System    - Time spent in system mode.
//   Idle      - Time spent in the idle task.
//   Iowait    - Time spent waiting for I/O to complete (since 2.5.41).
//   Irq       - Time servicing interrupts (since 2.6.0-test4).
//   Softirq   - Time servicing softirqs (since 2.6.0-test4).
//   Steal     - Stolen time, which is the time spent in other operating
//               systems when running a virtualized environment (since 2.6.11).
//   Guest     - Time spent running a virtual Cpu for guest operating
//               systems under the control of the Linux kernel (since 2.6.24).
//   GuestNice - Time spent running a niced guest (virtual Cpu for guest
//               operating systems under the control of the Linux kernel)
//               (since 2.6.33).
//   Total     - Total time.
// * Note: CPU time is measured in units of USER_HZ (1/100ths of a second on
//         (most architectures)
type CpuRawStats map[string]uint64

// CpuAvgStats represents *one* CPU statistics of a linux system.
// Keys:
//   User      - % of CPU time spent in user mode.
//   Nice      - % of CPU time spent in user mode with low priority (nice).
//   System    - % of CPU time spent in system mode.
//   Idle      - % of CPU time spent in the idle task.
//   Iowait    - % of CPU time spent waiting for I/O to complete (since 2.5.41).
//   Irq       - % of CPU servicing interrupts (since 2.6.0-test4).
//   Softirq   - % of CPU servicing softirqs (since 2.6.0-test4).
//   Steal     - % of stolen CPU time, which is the time spent in other operating
//               systems when running a virtualized environment (since 2.6.11).
//   Guest     - % of CPU time spent running a virtual Cpu for guest operating
//               systems under the control of the Linux kernel (since 2.6.24).
//   GuestNice - % of CPU time spent running a niced guest (virtual Cpu for guest
//               operating systems under the control of the Linux kernel)
//               (since 2.6.33).
//   Total     - Total time.
type CpuAvgStats map[string]float64

// CpusRawStats represents *all* the CPU raw statistics of a linux system.
// Keys:
//   Name - Name of the Cpu (as it is on /proc/stat: cpu, cpu0,...).
type CpusRawStats map[string]CpuRawStats

// CpusAvgStats represents *all* the CPU statistics of a linux system.
// Keys:
//   Name - Name of the Cpu (as it is on /proc/stat: cpu, cpu0,...).
type CpusAvgStats map[string]CpuAvgStats

// getCpuRawStats gets the Cpu raw stats of a linux system from the
// file /proc/stat
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
		cpuName, rawStats, err := parseCpuRawStats(stats)
		if err != nil {
			return nil, err
		}
		cpusRawStats[cpuName] = rawStats
	}

	return cpusRawStats, nil
}

// parseCpuRawStats parses the CPU stats as they are in the file /proc/stat.
// The stat file has the following format:
//   cpu  294 0 309 10612 71 30 0 0 0 0
// It returns:
// - cpuName is the name of the CPU (cpu, cpu0, cpu1, etc)
// - rawStats has the following format:
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
		rawStats[`total`] += stat
		switch i {
		case 1:
			rawStats[`user`] = stat
		case 2:
			rawStats[`nice`] = stat
		case 3:
			rawStats[`system`] = stat
		case 4:
			rawStats[`idle`] = stat
		case 5:
			rawStats[`iowait`] = stat
		case 6:
			rawStats[`irq`] = stat
		case 7:
			rawStats[`softirq`] = stat
		case 8:
			rawStats[`steal`] = stat
		case 9:
			rawStats[`guest`] = stat
		case 10:
			rawStats[`guestnice`] = stat
		}
	}

	return cpuName, rawStats, nil
}

// getCpuAvgStats calculates average between 2 CpusRawStats samples and returns
// the % CPU usage
func getCpuAvgStats(firstSample CpusRawStats, secondSample CpusRawStats) (cpusAvgStats CpusAvgStats, err error) {
	cpusAvgStats = CpusAvgStats{}

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

		cpuStats := CpuAvgStats{}
		timeDelta := float64(secondRawStats[`total`] - firstRawStats[`total`])
		// Calculate average between the two samples
		for key, secondValue := range secondRawStats {
			// Don't calculate average if the key is 'Total'
			if key == `Total` {
				continue
			}
			avg := float64(secondValue-firstRawStats[key]) * 100.00 / timeDelta
			avgStr := fmt.Sprintf("%3.2f", avg)
			cpuStats[key], err = strconv.ParseFloat(avgStr, 64)
			if err != nil {
				return nil, err
			}

		}
		cpuTotal := 100.00 - cpuStats[`idle`]
		cpuTotalStr := fmt.Sprintf("%3.2f", cpuTotal)
		cpuStats[`total`], err = strconv.ParseFloat(cpuTotalStr, 64)
		if err != nil {
			return nil, err
		}

		cpusAvgStats[cpuName] = cpuStats
	}

	return cpusAvgStats, nil
}

// getCpuStatsInterval returns the % CPU utilization between 2 samples.
// Time interval between the 2 samples is given in seconds.
func getCpuStatsInterval(interval int64) (cpusAvgStats CpusAvgStats, err error) {
	firstSample, err := getCpuRawStats()
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Duration(interval) * time.Second)

	secondSample, err := getCpuRawStats()
	if err != nil {
		return nil, err
	}

	cpusAvgStats, err = getCpuAvgStats(firstSample, secondSample)
	if err != nil {
		return nil, err
	}

	return cpusAvgStats, nil
}
