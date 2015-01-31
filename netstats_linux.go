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

// IfaceRawStats represents *one* network interface raw statistics of a
// linux system.
//
// Map keys:
//   rxbytes -  # of bytes.
//   rxpkts  -  # of packets.
//   rxerrs  -  # of errors that happend while receiving packets.
//   rxdrop  -  # of packets that were dropped.
//   rxfifo  -  # of FIFO overruns that happend on received packets.
//   rxframe -  # of carrier errors that happend on received packets.
//   rxcompr -  # of compressed packets received.
//   rxmulti -  # of multicast packets received.
//   txbytes -  # of bytes transmitted.
//   txpkts  -  # of packets transmitted.
//   txerrs  -  # of errors that happend while transmitting packets.
//   txdrop  -  # of packets that were dropped.
//   txfifo  -  # of FIFO overruns that happend on transmitted packets.
//   txcolls -  # of collisions that were detected.
//   txcarr  -  # of carrier errors that happend on transmitted packets.
//   txcompr -  # of compressed packets transmitted.
type IfaceRawStats map[string]uint64

// IfaceAvgStats represents *one* network interface statistics of a linux system.
//
// Map keys:
//   rxbytes -  # of bytes per second.
//   rxpkts  -  # of packets per second.
//   rxerrs  -  # of errors that happend while receiving packets per second.
//   rxdrop  -  # of packets that were dropped per second.
//   rxfifo  -  # of FIFO overruns that happend on received packets per second.
//   rxframe -  # of carrier errors that happend on received packets per second.
//   rxcompr -  # of compressed packets received per second.
//   rxmulti -  # of multicast packets received per second.
//   txbytes -  # of bytes transmitted per second.
//   txpkts  -  # of packets transmitted per second.
//   txerrs  -  # of errors that happend while transmitting packets per second.
//   txdrop  -  # of packets that were dropped per second.
//   txfifo  -  # of FIFO overruns that happend on transmitted packets per second.
//   txcolls -  # of collisions that were detected per second.
//   txcarr  -  # of carrier errors that happend on transmitted packets per second.
//   txcompr -  # of compressed packets transmitted per second.
type IfaceAvgStats map[string]float64

// NetRawStats represents *all* the network interfaces raw statistics of a linux system.
//
// Map keys:
//   Name - name of the network interface
type NetRawStats map[string]IfaceRawStats

// NetAvgStats represents *all* the network interfaces statistics of a linux system.
//
// Map keys:
//   Name - name of the network interface
type NetAvgStats map[string]IfaceAvgStats

// getNetRawStats gets the network interfaces raw statistics of a linux system from the
// file /proc/net/dev
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

// parseIfaceRawStats parses the network stats as they are in the file /proc/net/dev.
// It has the follogin format:
//  eth0:  178331 2395 0 0 0 0 0 0 257286 1876 0 0 0 0 0 0
//    lo:  166927  259 0 0 0 0 0 0 166927  259 0 0 0 0 0 0
// It returns:
//   - ifaceName, that is the name of the interface (lo, eth0,...)
//   - rawStats with the following format:
//       map[eth0:map[rxbytes:120 rxcompr:0 txdrop:0 rxpkts:2 rxerrs:0 txfifo:0
//                    rxdrop:0 rxframe:0 rxmulti:0 txbytes:276 txcolls:0 txcompr:0
//                    rxfifo:0 txpkts:2 txerrs:0 txcarr:0]
//             lo:map[rxpkts:0 rxerrs:0 txfifo:0 rxdrop:0 rxframe:0 rxmulti:0
//                    txbytes:0 txcolls:0 txcompr:0 rxfifo:0 txpkts:0 txerrs:0
//                    txcarr:0 rxbytes:0 rxcompr:0 txdrop:0]
//          ]
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

// getNetAvgStats calculates the network traffic average between 2 NetRawStats samples
func getNetAvgStats(firstSample NetRawStats, secondSample NetRawStats) (netAvgStats NetAvgStats, err error) {
	netAvgStats = NetAvgStats{}
	for ifaceName, secondRawStats := range secondSample {
		firstRawStats, ok := firstSample[ifaceName]
		if !ok {
			return nil, errors.New("The key " + ifaceName + " doesn't exist in the first sample of NetRawStats")
		}

		ifaceAvgStats := IfaceAvgStats{}
		timeDelta := float64(secondRawStats[`time`] - firstRawStats[`time`])
		for key, secondValue := range secondRawStats {
			if key == `time` {
				continue
			}
			avg := float64(secondValue-firstRawStats[key]) / timeDelta
			ifaceAvgStats[key] = avg
		}
		netAvgStats[ifaceName] = ifaceAvgStats
	}

	return netAvgStats, nil
}

// getNetAvgStatsInterval returns the network traffic average between 2 samples.
// Time interval between the 2 samples is given in seconds.
func getNetStatsInterval(interval int64) (netAvgStats NetAvgStats, err error) {
	firstSample, err := getNetRawStats()
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Duration(interval) * time.Second)

	secondSample, err := getNetRawStats()
	if err != nil {
		return nil, err
	}

	netAvgStats, err = getNetAvgStats(firstSample, secondSample)
	if err != nil {
		return nil, err
	}

	return netAvgStats, nil
}
