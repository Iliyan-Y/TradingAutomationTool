package strategies

import (
	helpers "stockApp/Helpers"
)

type BasicBuy struct {
	OpeningPrice float64
	ClosingPrice float64
	Threshold float64
}

type BasicSell struct {
	OpeningPrice float64
	ClosingPrice float64
	Threshold float64
}

func (s *BasicBuy) Validate() bool {
	return (helpers.PercentageDifference(s.OpeningPrice, s.ClosingPrice) <= s.Threshold)  
}

func (s *BasicSell) Validate() bool {
	return (helpers.PercentageDifference(s.OpeningPrice, s.ClosingPrice) >= s.Threshold)  
}