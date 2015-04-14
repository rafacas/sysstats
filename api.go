// Package sysstats provides system statistics.
package sysstats

// GetLoadAvg returns the load average of the system.
func GetLoadAvg() (LoadAvg, error) {
	return getLoadAvg()
}

// GetMemStats returns the memory statistics of the system.
func GetMemStats() (MemStats, error) {
	return getMemStats()
}

// GetCpuRawStats returns the CPUs statistics for the system at the moment
// the function is called.
func GetCpuRawStats() (CpusRawStats, error) {
	return getCpuRawStats()
}

// GetCpuAvgStats calculates average between 2 CPUs statistics samples and
// returns the % CPU usage
func GetCpuAvgStats(firstSample CpusRawStats, secondSample CpusRawStats) (CpusAvgStats, error) {
	return getCpuAvgStats(firstSample, secondSample)
}

// GetCpuStatsInterval returns the % CPU utilization between 2 samples where
// the sample interval is passed as an argument (in seconds).
func GetCpuStatsInterval(interval int64) (CpusAvgStats, error) {
	return getCpuStatsInterval(interval)
}

// GetNetRawStats returns all the network interfaces statistics of the system
func GetNetRawStats() (NetRawStats, error) {
	return getNetRawStats()
}

// GetNetAvgStats calculates average between 2 network stats samples
// and return the network traffic between them.
func GetNetAvgStats(firstSample NetRawStats, secondSample NetRawStats) (NetAvgStats, error) {
	return getNetAvgStats(firstSample, secondSample)
}

// GetNetStatsInterval returns the network traffic between 2 samples where the
// sample interval is passed as an argument (in seconds).
func GetNetStatsInterval(interval int64) (NetAvgStats, error) {
	return getNetStatsInterval(interval)
}

// GetDiskUsage gets an array (one element per partition) with the disk
// usage of the system
func GetDiskUsage() ([]DiskUsage, error) {
	return getDiskUsage()
}

// GetDiskRawStats gets the disk IO stats of the system at the moment
// the function is called.
func GetDiskRawStats() ([]DiskRawStats, error) {
	return getDiskRawStats()
}

// GetDiskAvgStats calculates the average between 2 DiskRawStats samples and
// returns the number of IOs per second.
func GetDiskAvgStats(firstSampleArr []DiskRawStats, secondSampleArr []DiskRawStats) ([]DiskAvgStats, error) {
	return getDiskAvgStats(firstSampleArr, secondSampleArr)
}

// GetDiskStatsInterval returns the IO average between 2 samples where
// the sample interval is passed as an argument (in seconds).
func GetDiskStatsInterval(interval int64) ([]DiskAvgStats, error) {
	return getDiskStatsInterval(interval)
}

// GetSockStats returns the socket statistics of the system.
func GetSockStats() (SockStats, error) {
	return getSockStats()
}

// GetSysInfo returns the system info (as hostname, OS type, etc).
func GetSysInfo() (SysInfo, error) {
	return getSysInfo()
}

// GetFileStats returns the file statistics of the system.
func GetFileStats() (FileStats, error) {
	return getFileStats()
}

// GetProcRawStats returns the processes stats of the system.
func GetProcRawStats() (ProcRawStats, error) {
	return getProcRawStats()
}

// GetProcAvgStats calculates the average between 2 processes stats samples.
func GetProcAvgStats(firstSample ProcRawStats, secondSample ProcRawStats) (ProcAvgStats, error) {
	return getProcAvgStats(firstSample, secondSample)
}

// GetProcStatsInterval returns the processes stats average between 2 samples
// where the sample interval is passed as an argument (in seconds).
func GetProcStatsInterval(interval int64) (ProcAvgStats, error) {
	return getProcStatsInterval(interval)
}
