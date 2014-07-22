package sysstats

func GetLoadAvg() (*LoadAvg, error) {
	return getLoadAvg()
}
