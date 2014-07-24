// +build darwin

package sysstats

import (
	"errors"
	"runtime"
)

// MemStat represents the memory statistics on an OSX system
type MemStats map[string]uint64

// getMemStats gets the memory stats of an OSX system
func getMemStats() (memStats MemStats, err error) {
	return nil, errors.New("getMemStats: " + runtime.GOOS + " not supported yet")
}
