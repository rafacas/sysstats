package sysstats

func GetLoadAvg() (LoadAvg, error) {
	return getLoadAvg()
}

func GetMemStats() (MemStats, error) {
	return getMemStats()
}

func GetCpuStats() (CpusRawStats, error) {
	return getCpuStats()
}
