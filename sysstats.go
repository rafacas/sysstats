package sysstats

func GetLoadAvg() (LoadAvg, error) {
	return getLoadAvg()
}

func GetMemStats() (MemStats, error) {
	return getMemStats()
}

func GetCpuRawStats() (CpusRawStats, error) {
	return getCpuRawStats()
}

func GetCpuStats(firstSample CpusRawStats, secondSample CpusRawStats) (CpusStats, error) {
	return getCpuStats(firstSample, secondSample)
}

func GetCpuStatsInterval(interval int64) (CpusStats, error) {
	return getCpuStatsInterval(interval)
}

func GetNetRawStats() (NetRawStats, error) {
	return getNetRawStats()
}

func GetNetStats(firstSample NetRawStats, secondSample NetRawStats) (NetStats, error) {
	return getNetStats(firstSample, secondSample)
}

func GetNetStatsInterval(interval int64) (NetStats, error) {
	return getNetStatsInterval(interval)
}

func GetDiskUsage() ([]DiskUsage, error) {
	return getDiskUsage()
}

func GetDiskRawStats() ([]DiskRawStats, error) {
	return getDiskRawStats()
}

func GetDiskStats(firstSample DiskRawStats, secondSample DiskRawStats) (DiskStats, error) {
	return getDiskStats(firstSample, secondSample)
}

func GetDiskStatsInterval(interval int64) ([]DiskStats, error) {
	return getDiskStatsInterval(interval)
}

func GetSockStats() (SockStats, error) {
	return getSockStats()
}

func GetSysInfo() (SysInfo, error) {
	return getSysInfo()
}
