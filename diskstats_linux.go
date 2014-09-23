// +build linux

package sysstats

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// DiskRawStats represents the disk IO raw statistics of a linux system.
type DiskRawStats struct {
	Major        int    // Major number for the disk
	Minor        int    // Minor number for the disk
	Name         string // Disk name
	ReadIOs      uint64 // # of reads completed since boot
	ReadMerges   uint64 // # of reads merged since boot
	ReadSectors  uint64 // # of sectors read since boot
	ReadTicks    uint64 // # of milliseconds spent reading since boot
	WriteIOs     uint64 // # of writes completed since boot
	WriteMerges  uint64 // # of writes merged since boot
	WriteSectors uint64 // # of sectors written since boot
	WriteTicks   uint64 // # of milliseconds spent writing since boot
	InFlight     uint64 // # of I/Os currently in progress
	IOTicks      uint64 // # of milliseconds spent doing I/Os since boot
	TimeInQueue  uint64 // Weighted # of milliseconds spent doing I/Os since boot
	SampleTime   int64  // Time when the sample was taken
}

// DiskAvgStats represents the average disk IO statistics (per second) of a
// linux system.
type DiskAvgStats struct {
	Major       int     // Major number for the disk
	Minor       int     // Minor number for the disk
	Name        string  // Disk name
	ReadIOs     float64 // # of reads completed per second
	ReadMerges  float64 // # of reads merged per second
	ReadBytes   float64 // # of bytes read per second
	WriteIOs    float64 // # of writes completed per second
	WriteMerges float64 // # of writes merged per second
	WriteBytes  float64 // # of bytes written per second
	InFlight    uint64  // # of I/Os currently in progress
	IOTicks     uint64  // # of milliseconds spent doing I/Os
	TimeInQueue uint64  // Weighted # of milliseconds spent doing I/Os
}

// getDiskRawStats gets the disk IO stats of a linux system from the
// file /proc/diskstats
func getDiskRawStats() (diskRawStatsArr []DiskRawStats, err error) {
	file, err := os.Open("/proc/diskstats")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	diskRawStatsArr = make([]DiskRawStats, 0, 5)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	now := time.Now().Unix()
	for scanner.Scan() {
		line := scanner.Text()
		diskRawStats, err := parseDiskRawStats(line)
		if err != nil {
			return diskRawStatsArr, err
		}
		diskRawStats.SampleTime = now
		diskRawStatsArr = append(diskRawStatsArr, diskRawStats)
	}

	return diskRawStatsArr, nil
}

// parseDiskRawStats parses the disk stats.
// The file /proc/diskstats has the following format:
//   7       0 loop0 0 0 0 0 0 0 0 0 0 0 0
//   7       1 loop1 0 0 0 0 0 0 0 0 0 0 0
//   7       2 loop2 0 0 0 0 0 0 0 0 0 0 0
//   7       3 loop3 0 0 0 0 0 0 0 0 0 0 0
//   7       4 loop4 0 0 0 0 0 0 0 0 0 0 0
//   7       5 loop5 0 0 0 0 0 0 0 0 0 0 0
//   7       6 loop6 0 0 0 0 0 0 0 0 0 0 0
//   7       7 loop7 0 0 0 0 0 0 0 0 0 0 0
//   8       0 sda 4222 4373 293854 48992 676 1024 13428 2016 0 1744 51004
//   8       1 sda1 287 322 2296 68 6 0 12 0 0 68 68
//   8       2 sda2 2 0 4 0 0 0 0 0 0 0 0
//   8       5 sda5 3748 4051 290074 48904 587 1024 13416 2016 0 1676 50916
// 252       0 dm-0 7516 0 287642 65724 1613 0 13416 4212 0 1644 69936
// 252       1 dm-1 224 0 1792 28 0 0 0 0 0 28 28
func parseDiskRawStats(stats string) (diskRawStats DiskRawStats, err error) {
	diskRawStats = DiskRawStats{}

	fields := strings.Fields(stats)

	// Check there are 14 fields
	if len(fields) != 14 {
		return diskRawStats, errors.New("Couldn't parse disk stats because there aren't 14 fields")
	}

	// Parse fields
	for i := 0; i < len(fields); i++ {
		field := fields[i]
		switch i {
		case 0:
			major, _ := strconv.ParseInt(field, 10, strconv.IntSize)
			diskRawStats.Major = int(major)
		case 1:
			minor, _ := strconv.ParseInt(field, 10, strconv.IntSize)
			diskRawStats.Minor = int(minor)
		case 2:
			diskRawStats.Name = fields[2]
		case 3:
			readIOs, _ := strconv.ParseUint(field, 10, 64)
			diskRawStats.ReadIOs = readIOs
		case 4:
			readMerges, _ := strconv.ParseUint(field, 10, 64)
			diskRawStats.ReadMerges = readMerges
		case 5:
			readSectors, _ := strconv.ParseUint(field, 10, 64)
			diskRawStats.ReadSectors = readSectors
		case 6:
			readTicks, _ := strconv.ParseUint(field, 10, 64)
			diskRawStats.ReadTicks = readTicks
		case 7:
			writeIOs, _ := strconv.ParseUint(field, 10, 64)
			diskRawStats.WriteIOs = writeIOs
		case 8:
			writeMerges, _ := strconv.ParseUint(field, 10, 64)
			diskRawStats.WriteMerges = writeMerges
		case 9:
			writeSectors, _ := strconv.ParseUint(field, 10, 64)
			diskRawStats.WriteSectors = writeSectors
		case 10:
			writeTicks, _ := strconv.ParseUint(field, 10, 64)
			diskRawStats.WriteTicks = writeTicks
		case 11:
			inFlight, _ := strconv.ParseUint(field, 10, 64)
			diskRawStats.InFlight = inFlight
		case 12:
			ioTicks, _ := strconv.ParseUint(field, 10, 64)
			diskRawStats.IOTicks = ioTicks
		case 13:
			timeInQueue, _ := strconv.ParseUint(field, 10, 64)
			diskRawStats.TimeInQueue = timeInQueue
		}
	}

	return diskRawStats, nil
}

// getDiskAvgStats calculates the average between 2 DiskRawStats samples and returns
// a DiskAvgStats variable with the number of IOs per second.
func getDiskAvgStats(firstSample DiskRawStats, secondSample DiskRawStats) (diskAvgStats DiskAvgStats, err error) {
	diskAvgStats = DiskAvgStats{}

	timeDelta := float64(secondSample.SampleTime - firstSample.SampleTime)

	// Check the samples are from the same disk
	if firstSample.Major != secondSample.Major ||
		firstSample.Minor != secondSample.Minor ||
		firstSample.Name != secondSample.Name {
		msg := fmt.Sprintf("The samples are from different disks: \n\tfirstSample -> %d %d %s \n\t"+
			"secondSample -> %d %d %s\n", firstSample.Major, firstSample.Minor, firstSample.Name,
			secondSample.Major, secondSample.Minor, secondSample.Name)
		return DiskAvgStats{}, errors.New(msg)
	} else {
		diskAvgStats.Major = firstSample.Major
		diskAvgStats.Minor = firstSample.Minor
		diskAvgStats.Name = firstSample.Name
	}

	// Calculate average between the 2 samples
	diskAvgStats.ReadIOs = float64(secondSample.ReadIOs-firstSample.ReadIOs) / timeDelta
	diskAvgStats.ReadMerges = float64(secondSample.ReadMerges-firstSample.ReadMerges) / timeDelta
	diskAvgStats.ReadBytes = float64((secondSample.ReadSectors*512)-(firstSample.ReadSectors*512)) / timeDelta
	diskAvgStats.WriteIOs = float64(secondSample.WriteIOs-firstSample.WriteIOs) / timeDelta
	diskAvgStats.WriteMerges = float64(secondSample.WriteMerges-firstSample.WriteMerges) / timeDelta
	diskAvgStats.WriteBytes = float64((secondSample.WriteSectors*512)-(firstSample.WriteSectors*512)) / timeDelta

	diskAvgStats.InFlight = secondSample.InFlight
	diskAvgStats.TimeInQueue = secondSample.TimeInQueue - firstSample.TimeInQueue

	return diskAvgStats, nil
}

// getDiskStatsInterval returns the IO average between 2 samples.
// Time interval between the 2 samples is given in seconds.
func getDiskStatsInterval(interval int64) (diskAvgStatsArr []DiskAvgStats, err error) {
	firstSampleArr, err := getDiskRawStats()
	if err != nil {
		return nil, err
	}

	diskAvgStatsArr = make([]DiskAvgStats, 0, len(firstSampleArr)-1)

	time.Sleep(time.Duration(interval) * time.Second)

	secondSampleArr, err := getDiskRawStats()
	if err != nil {
		return nil, err
	}

	for _, firstSample := range firstSampleArr {
		diskName := firstSample.Name
		for _, secondSample := range secondSampleArr {
			if secondSample.Name == diskName {
				diskAvgStats, err := getDiskAvgStats(firstSample, secondSample)
				if err != nil {
					return nil, err
				}
				diskAvgStatsArr = append(diskAvgStatsArr, diskAvgStats)
				break
			} else {
				continue
			}
		}
	}

	return diskAvgStatsArr, nil
}
