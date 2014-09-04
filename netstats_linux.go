// +build linux

package sysstats

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// IfaceRawStats represents *one* network interface raw statistics of a linux system.
type IfaceRawStats map[string]uint64

// IfaceStats represents *one* network interface statistics (%) of a linux system.
type IfaceStats map[string]float64

// NetRawStats represents *all* the network interfaces raw statistics of a linux system.
type NetRawStats map[string]IfaceRawStats

// NetStats represents *all* the network interfaces statistics (%) of a linux system.
type NetStats map[string]IfaceStats

// getNetRawStats gets the network interfaces raw statistics of a linux system from the
// file /proc/net/dev
// It returns a NetRawStats var. It is a map which keys are the 'kernel name'
// of the network interfaces (lo, eth0, eth1, etc). The value of each key is a
// IfaceRawStats var with the statistics for that interface.
// NetRawStats has the following key:
//  Name    - name of the network interface
// IfaceRawStats has the following keys:
// rxbytes  -  Number of bytes.
// rxpkts   -  Number of packets.
// rxerrs   -  Number of errors that happend while receiving packets.
// rxdrop   -  Number of packets that were dropped.
// rxfifo   -  Number of FIFO overruns that happend on received packets.
// rxframe  -  Number of carrier errors that happend on received packets.
// rxcompr  -  Number of compressed packets received.
// rxmulti  -  Number of multicast packets received.
// txbytes  -  Number of bytes transmitted.
// txpkts   -  Number of packets transmitted.
// txerrs   -  Number of errors that happend while transmitting packets.
// txdrop   -  Number of packets that were dropped.
// txfifo   -  Number of FIFO overruns that happend on transmitted packets.
// txcolls  -  Number of collisions that were detected.
// txcarr   -  Number of carrier errors that happend on transmitted packets.
// txcompr  -  Number of compressed packets transmitted.
func getNetRawStats() (netRawStats NetRawStats, err error) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	netRawStats = NetRawStats{}

	re := regexp.MustCompile(`^\s*(.+?):\s*(.*)`)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	now := time.Now().Unix()
	for scanner.Scan() {
		line := scanner.Text()
		stats := re.FindString(line)
		if stats == "" {
			// No match
			continue
		}
		ifaceName, rawStats, err := parseIfaceRawStats(stats)
		if err != nil {
			return nil, err
		}
		rawStats[`time`] = uint64(now)
		netRawStats[ifaceName] = rawStats
	}

	return netRawStats, nil
}

// parseIfaceRawStats parses the network stats as they are in the file /proc/net/dev
// It has the follogin format:
//  eth0:  178331 2395 0 0 0 0 0 0 257286 1876 0 0 0 0 0 0
//    lo:  166927  259 0 0 0 0 0 0 166927  259 0 0 0 0 0 0
// It returns the ifaceName, that is the name of the interface (lo, eth0,...)
// and the rawStats with the following format:
//   map[eth0:map[rxbytes:120 rxcompr:0 txdrop:0 rxpkts:2 rxerrs:0 txfifo:0
//                rxdrop:0 rxframe:0 rxmulti:0 txbytes:276 txcolls:0 txcompr:0
//                rxfifo:0 txpkts:2 txerrs:0 txcarr:0]
//         lo:map[rxpkts:0 rxerrs:0 txfifo:0 rxdrop:0 rxframe:0 rxmulti:0
//                txbytes:0 txcolls:0 txcompr:0 rxfifo:0 txpkts:0 txerrs:0
//                txcarr:0 rxbytes:0 rxcompr:0 txdrop:0]
//      ]
func parseIfaceRawStats(stats string) (ifaceName string, rawStats IfaceRawStats,
	err error) {

	rawStats = IfaceRawStats{}

	fields := strings.Fields(stats)
	ifaceName = fields[0]
	// Trim the trailing ':'
	if last := len(ifaceName) - 1; last >= 0 && ifaceName[last] == ':' {
		ifaceName = ifaceName[:last]
	}

	for i := 1; i < len(fields); i++ {
		stat, err := strconv.ParseUint(fields[i], 10, 64)
		if err != nil {
			return "", nil, err
		}

		switch i {
		case 1:
			rawStats[`rxbytes`] = stat
		case 2:
			rawStats[`rxpkts`] = stat
		case 3:
			rawStats[`rxerrs`] = stat
		case 4:
			rawStats[`rxdrop`] = stat
		case 5:
			rawStats[`rxfifo`] = stat
		case 6:
			rawStats[`rxframe`] = stat
		case 7:
			rawStats[`rxcompr`] = stat
		case 8:
			rawStats[`rxmulti`] = stat
		case 9:
			rawStats[`txbytes`] = stat
		case 10:
			rawStats[`txpkts`] = stat
		case 11:
			rawStats[`txerrs`] = stat
		case 12:
			rawStats[`txdrop`] = stat
		case 13:
			rawStats[`txfifo`] = stat
		case 14:
			rawStats[`txcolls`] = stat
		case 15:
			rawStats[`txcarr`] = stat
		case 16:
			rawStats[`txcompr`] = stat
		}
	}

	return ifaceName, rawStats, nil
}

// getNetStats calculates the network traffic average between 2 NetRawStats samples
func getNetStats(firstSample NetRawStats, secondSample NetRawStats) (netStats NetStats, err error) {
	netStats = NetStats{}
	for ifaceName, secondRawStats := range secondSample {
		firstRawStats, ok := firstSample[ifaceName]
		if !ok {
			return nil, errors.New("The key " + ifaceName + " doesn't exist in the first sample of NetRawStats")
		}

		ifaceStats := IfaceStats{}
		timeDelta := float64(secondRawStats[`time`] - firstRawStats[`time`])
		for key, secondValue := range secondRawStats {
			if key == `time` {
				continue
			}
			avg := float64(secondValue-firstRawStats[key]) / timeDelta
			ifaceStats[key] = avg
		}
		netStats[ifaceName] = ifaceStats
	}

	return netStats, nil
}

// getNetStatsInterval returns the network traffic average between 2 samples taken in a
// time interval (given in seconds)
func getNetStatsInterval(interval int64) (netStats NetStats, err error) {
	firstSample, err := getNetRawStats()
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Duration(interval) * time.Second)

	secondSample, err := getNetRawStats()
	if err != nil {
		return nil, err
	}

	netStats, err = getNetStats(firstSample, secondSample)
	if err != nil {
		return nil, err
	}

	return netStats, nil
}
