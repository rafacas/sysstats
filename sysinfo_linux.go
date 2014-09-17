// +build linux

package sysstats

import (
	"io/ioutil"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// SysInfo represents the linux system info
type SysInfo struct {
	Hostname  string
	Domain    string
	OsType    string
	OsRelease string
	OsVersion string
	OsArch    string
	Uptime    float64
}

// getSysInfo gets the system info
func getSysInfo() (sysInfo SysInfo, err error) {
	sysInfo = SysInfo{}

	// Hostname
	hostname, err := getHostname()
	if err != nil {
		return SysInfo{}, err
	}
	sysInfo.Hostname = hostname

	// Domain
	domain, err := getDomain()
	if err != nil {
		return SysInfo{}, err
	}
	sysInfo.Domain = domain

	// OS type
	osType, err := getOsType()
	if err != nil {
		return SysInfo{}, err
	}
	sysInfo.OsType = osType

	// OS relase
	osRelease, err := getOsRelease()
	if err != nil {
		return SysInfo{}, err
	}
	sysInfo.OsRelease = osRelease

	// OS version
	osVersion, err := getOsVersion()
	if err != nil {
		return SysInfo{}, err
	}
	sysInfo.OsVersion = osVersion

	// OS arch
	osArch, err := getOsArch()
	if err != nil {
		return SysInfo{}, err
	}
	sysInfo.OsArch = osArch

	// Uptime
	uptime, err := getUptime()
	if err != nil {
		return SysInfo{}, err
	}
	sysInfo.Uptime = uptime

	return sysInfo, nil
}

func getHostname() (hostname string, err error) {
	content, err := ioutil.ReadFile("/proc/sys/kernel/hostname")
	if err != nil {
		return "", err
	}

	hostname = strings.TrimSpace(string(content))
	return hostname, nil
}

func getDomain() (domain string, err error) {
	content, err := ioutil.ReadFile("/proc/sys/kernel/domainname")
	if err != nil {
		return "", err
	}

	domain = strings.TrimSpace(string(content))
	return domain, nil
}

func getOsType() (osType string, err error) {
	content, err := ioutil.ReadFile("/proc/sys/kernel/ostype")
	if err != nil {
		return "", err
	}

	osType = strings.TrimSpace(string(content))
	return osType, nil
}

func getOsRelease() (osRelease string, err error) {
	content, err := ioutil.ReadFile("/proc/sys/kernel/osrelease")
	if err != nil {
		return "", err
	}

	osRelease = strings.TrimSpace(string(content))
	return osRelease, nil
}

func getOsVersion() (osVersion string, err error) {
	content, err := ioutil.ReadFile("/proc/sys/kernel/version")
	if err != nil {
		return "", err
	}

	osVersion = strings.TrimSpace(string(content))
	return osVersion, nil
}

func getOsArch() (osArch string, err error) {
	// Check `uname` path
	uname, err := exec.LookPath("uname")
	if err != nil {
		return "", err
	}

	// Run `uname -m` to get the OS architecture
	out, err := exec.Command(uname, "-m").Output()
	if err != nil {
		return "", err
	}

	osArch = strings.TrimSpace(string(out))
	return osArch, nil
}

func getUptime() (uptime float64, err error) {
	content, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		return -1, err
	}

	re := regexp.MustCompile(`^(\d*\.?\d*)\s+`)
	match := re.FindStringSubmatch(string(content))
	if err != nil {
		return -1, err
	}

	uptime, err = strconv.ParseFloat(match[1], 64)
	if err != nil {
		return -1, err
	}

	return uptime, nil
}
