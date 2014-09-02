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
