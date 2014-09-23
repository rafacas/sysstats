// +build linux

package sysstats

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
)

// SockStats represents the socket statistics of a linux system.
type SockStats struct {
	Used        uint64 // Total number of used sockets
	TcpInUse    uint64 // TCP sockets in use
	TcpOrphaned uint64 // TCP sockets orphaned
	TcpTimeWait uint64 // TCP sockets in TIME_WAIT
	UdpInUse    uint64 // UDP sockets in use
	Raw         uint64 // RAW sockets in use
	IpFrag      uint64 // # of IP fragments in use
}

// getSockStats gets the socket statistics of a linux system from the file
// /proc/net/sockstat
func getSockStats() (sockStats SockStats, err error) {
	file, err := os.Open("/proc/net/sockstat")
	if err != nil {
		return SockStats{}, err
	}
	defer file.Close()

	sockStats = SockStats{}
	reSock := regexp.MustCompile(`sockets:\s+used\s+(\d+)`)
	reTcp := regexp.MustCompile(`TCP:\s+inuse\s+(\d+)\s+orphan\s+(\d+)\s+tw\s+(\d+)`)
	reUdp := regexp.MustCompile(`UDP:\s+inuse\s+(\d+)`)
	reRaw := regexp.MustCompile(`RAW:\s+inuse\s+(\d+)`)
	reFrag := regexp.MustCompile(`FRAG:\s+inuse\s+(\d+)`)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if stat := reSock.FindStringSubmatch(line); stat != nil {
			sockUsed, err := strconv.ParseUint(stat[1], 10, 64)
			if err != nil {
				return SockStats{}, err
			}
			sockStats.Used = sockUsed
		} else if stat := reTcp.FindStringSubmatch(line); stat != nil {
			tcpInUse, err := strconv.ParseUint(stat[1], 10, 64)
			if err != nil {
				return SockStats{}, err
			}
			sockStats.TcpInUse = tcpInUse
			tcpOrphaned, err := strconv.ParseUint(stat[2], 10, 64)
			if err != nil {
				return SockStats{}, err
			}
			sockStats.TcpOrphaned = tcpOrphaned
			tcpTimeWait, err := strconv.ParseUint(stat[3], 10, 64)
			if err != nil {
				return SockStats{}, err
			}
			sockStats.TcpTimeWait = tcpTimeWait
		} else if stat := reUdp.FindStringSubmatch(line); stat != nil {
			udpInUse, err := strconv.ParseUint(stat[1], 10, 64)
			if err != nil {
				return SockStats{}, err
			}
			sockStats.UdpInUse = udpInUse
		} else if stat := reRaw.FindStringSubmatch(line); stat != nil {
			raw, err := strconv.ParseUint(stat[1], 10, 64)
			if err != nil {
				return SockStats{}, err
			}
			sockStats.Raw = raw
		} else if stat := reFrag.FindStringSubmatch(line); stat != nil {
			ipFrag, err := strconv.ParseUint(stat[1], 10, 64)
			if err != nil {
				return SockStats{}, err
			}
			sockStats.IpFrag = ipFrag
		}
	}

	return sockStats, nil
}
