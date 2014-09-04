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
