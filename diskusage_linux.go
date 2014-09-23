// build +linux

package sysstats

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

// DiskUsage represents a file system disk space usage
type DiskUsage struct {
	FileSystem string
	Type       string
	Total      uint64
	Used       uint64
	Available  uint64
	UsedPer    uint64
	MountedOn  string
}

// getDiskUsage gets the disk usage of a linux system running the command:
//   df -kTP
// where:
//     -k: block size = 1K
//     -T: prints the file system type
//     -P: uses the POSIX output format
// It returns an array of DiskUsage elements (as many elements as file systems
// has the OS)
func getDiskUsage() (diskUsageArr []DiskUsage, err error) {
	diskUsageArr = make([]DiskUsage, 0, 5)

	// Check df exists
	df, err := exec.LookPath("df")
	if err != nil {
		return diskUsageArr, err
	}

	// Run df -kTP
	out, err := exec.Command(df, "-kTP").Output()

	if err != nil {
		return diskUsageArr, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	scanner.Split(bufio.ScanLines)
	// Filter the header
	scanner.Scan()
	for scanner.Scan() {
		line := scanner.Text()
		diskUsage, err := parseDiskUsage(line)
		if err != nil {
			return diskUsageArr, err
		}

		diskUsageArr = append(diskUsageArr, diskUsage)
	}

	return diskUsageArr, nil
}

// parseDiskUsage parses the filesystem disk space usage reported by df.
// The format of the usage sent as argument has the following format:
//   Filesystem                      Type     1024-blocks      Used Available Capacity Mounted on
//   /dev/mapper/ubuntu1404--vg-root ext4         7359808   1866928   5095980      27% /
//   none                            tmpfs              4         0         4       0% /sys/fs/cgroup
//   udev                            devtmpfs      239692         4    239688       1% /dev
//   tmpfs                           tmpfs          50184       348     49836       1% /run
//   none                            tmpfs           5120         0      5120       0% /run/lock
//   none                            tmpfs         250916         0    250916       0% /run/shm
//   none                            tmpfs         102400         0    102400       0% /run/user
//   /dev/sda1                       ext2          240972     36441    192090      16% /boot
func parseDiskUsage(usage string) (diskUsage DiskUsage, err error) {
	diskUsage = DiskUsage{}

	fields := strings.Fields(usage)

	// Check there are 7 fields
	if len(fields) != 7 {
		return DiskUsage{}, errors.New("Couldn't parse disk usage because there aren't 7 fields")
	}

	// Parse fields
	for i := 0; i < len(fields); i++ {
		field := fields[i]
		switch i {
		case 0:
			diskUsage.FileSystem = field
		case 1:
			diskUsage.Type = field
		case 2:
			value, err := strconv.ParseUint(field, 10, 64)
			if err != nil {
				return DiskUsage{}, err
			}
			diskUsage.Total = value
		case 3:
			value, err := strconv.ParseUint(field, 10, 64)
			if err != nil {
				return DiskUsage{}, err
			}
			diskUsage.Used = value
		case 4:
			value, err := strconv.ParseUint(field, 10, 64)
			if err != nil {
				return DiskUsage{}, err
			}
			diskUsage.Available = value
		case 5:
			// Trim trailing '%'
			if last := len(field) - 1; last >= 0 && field[last] == '%' {
				field = field[:last]
			}
			value, err := strconv.ParseUint(field, 10, 64)
			if err != nil {
				return DiskUsage{}, err
			}
			diskUsage.UsedPer = value
		case 6:
			diskUsage.MountedOn = field
		}
	}

	return diskUsage, nil
}
