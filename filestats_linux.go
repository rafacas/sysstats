// +build linux

package sysstats

import (
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
)

// FileStats represents the file descriptor stats
type FileStats struct {
	FhAlloc uint64 // # of allocated file handlers (# files currently opened)
	FhFree  uint64 // # of free file handlers
	FhMax   uint64 // maximum # of file handlers
	InAlloc uint64 // # of inodes the system has allocated
	InFree  uint64 // # of free inodes
}

// getFileStats gets the file statistics of a linux system from the files:
// /proc/sys/fs/file-nr and /proc/sys/fs/inode-nr
func getFileStats() (fileStats FileStats, err error) {
	fileStats = FileStats{}

	// Get file handler stats
	content, err := ioutil.ReadFile("/proc/sys/fs/file-nr")
	if err != nil {
		return FileStats{}, err
	}

	fields := strings.Fields(strings.TrimSpace(string(content)))
	if len(fields) != 3 {
		return FileStats{}, errors.New("Error parsing file /proc/sys/fs/file-nr. It should have 3 fields")
	}
	fileStats.FhAlloc, err = strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return FileStats{}, err
	}
	fileStats.FhFree, err = strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return FileStats{}, err
	}
	fileStats.FhMax, err = strconv.ParseUint(fields[2], 10, 64)
	if err != nil {
		return FileStats{}, err
	}

	// Get the inode stats
	content, err = ioutil.ReadFile("/proc/sys/fs/inode-nr")
	if err != nil {
		return FileStats{}, err
	}

	fields = strings.Fields(strings.TrimSpace(string(content)))
	if len(fields) != 2 {
		return FileStats{}, errors.New("Error parsing file /proc/sys/fs/inode-nr. It should have 2 fields")
	}
	fileStats.InAlloc, err = strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return FileStats{}, err
	}
	fileStats.InFree, err = strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return FileStats{}, err
	}

	return fileStats, nil
}
