// +build linux

package sysstats

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

// MemStat represents the memory statistics on a linux system
// The following are the keys of the map:
// memused      -  Total size of used memory in kilobytes.
// memfree      -  Total size of free memory in kilobytes.
// memusedper   -  Total size of used memory in percent.
// memtotal     -  Total size of memory in kilobytes.
// buffers      -  Total size of buffers used from memory in kilobytes.
// cached       -  Total size of cached memory in kilobytes.
// realfree     -  Total size of memory is real free (memfree + buffers + cached).
// realfreeper  -  Total size of memory is real free in percent of total memory.
// swapused     -  Total size of swap space is used is kilobytes.
// swapfree     -  Total size of swap space is free in kilobytes.
// swapusedper  -  Total size of swap space is used in percent.
// swaptotal    -  Total size of swap space in kilobytes.
// swapcached   -  Memory that once was swapped out, is swapped back in but still
//                 also is in the swapfile.
// active       -  Memory that has been used more recently and usually not
//                 reclaimed unless absolutely necessary.
// inactive     -  Memory which has been less recently used and is more eligible
//                 to be reclaimed for other purposes.
//
// The following statistics are only available for kernels >= 2.6.
// slab         -  Total size of memory in kilobytes that used by kernel for data
//                 structure allocations.
// dirty        -  Total size of memory pages in kilobytes that waits to be
//                 written back to disk.
// mapped       -  Total size of memory in kilobytes that is mapped by devices or
//                 libraries with mmap.
// writeback    -  Total size of memory that was written back to disk.
// committed_as -  The amount of memory presently allocated on the system.
//
// The following statistic is only available for kernels >= 2.6.9.
// commitlimit  -  Total amount of memory currently available to be allocated on
//                 the system.
type MemStats map[string]uint64

// getMemStats gets the memory stats of a linux system from the
// file /proc/meminfo
func getMemStats() (memStats MemStats, err error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	memStats = MemStats{}
	re := regexp.MustCompile(`^((?:Mem|Swap)(?:Total|Free)|Buffers|Cached|` +
		`SwapCached|Active|Inactive|Dirty|Writeback|Mapped|Slab|` +
		`Commit(?:Limit|ted_AS)):\s*(\d+)`)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		stat := re.FindStringSubmatch(line)
		if stat == nil {
			// No match
			continue
		}
		key := stat[1]
		value, err := strconv.ParseUint(stat[2], 10, 64)
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			memStats[key] = value
		}
	}

	return memStats, nil
}
