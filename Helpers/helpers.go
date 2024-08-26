package helpers

import domain "stockApp/Domain"

func PercentageDifference(oldValue, newValue float64) float64 {
	if oldValue == 0 {
		return float64(0)
	}
	return ((newValue - oldValue) / oldValue) * 100
}

func FilterTradeByIndex(trades *[]domain.Trade, index string) *domain.Trade {
	for _, trade := range *trades {
			if trade.Index == index {
					return &trade
			}
	}
	return nil
}

func NotifyCashIn() {
	//todo
}