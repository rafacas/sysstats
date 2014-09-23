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

func GetCpuAvgStats(firstSample CpusRawStats, secondSample CpusRawStats) (CpusAvgStats, error) {
	return getCpuAvgStats(firstSample, secondSample)
}

func GetCpuStatsInterval(interval int64) (CpusAvgStats, error) {
	return getCpuStatsInterval(interval)
}

func GetNetRawStats() (NetRawStats, error) {
	return getNetRawStats()
}

func GetNetAvgStats(firstSample NetRawStats, secondSample NetRawStats) (NetAvgStats, error) {
	return getNetAvgStats(firstSample, secondSample)
}

func GetNetStatsInterval(interval int64) (NetAvgStats, error) {
	return getNetStatsInterval(interval)
}

func GetDiskUsage() ([]DiskUsage, error) {
	return getDiskUsage()
}

func GetDiskRawStats() ([]DiskRawStats, error) {
	return getDiskRawStats()
}

func GetDiskAvgStats(firstSample DiskRawStats, secondSample DiskRawStats) (DiskAvgStats, error) {
	return getDiskAvgStats(firstSample, secondSample)
}

func GetDiskStatsInterval(interval int64) ([]DiskAvgStats, error) {
	return getDiskStatsInterval(interval)
}

func GetSockStats() (SockStats, error) {
	return getSockStats()
}

func GetSysInfo() (SysInfo, error) {
	return getSysInfo()
}

func GetFileStats() (FileStats, error) {
	return getFileStats()
}

func GetProcRawStats() (ProcRawStats, error) {
	return getProcRawStats()
}

func GetProcAvgStats(firstSample ProcRawStats, secondSample ProcRawStats) (ProcAvgStats, error) {
	return getProcAvgStats(firstSample, secondSample)
}

func GetProcStatsInterval(interval int64) (ProcAvgStats, error) {
	return getProcStatsInterval(interval)
}
