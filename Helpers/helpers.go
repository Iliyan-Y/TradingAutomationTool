package helpers

func PercentageDifference(oldValue, newValue float64) float64 {
	if oldValue == 0 {
		return float64(0)
	}
	return ((newValue - oldValue) / oldValue) * 100
}

func NotifyCashIn() {
	//todo
}