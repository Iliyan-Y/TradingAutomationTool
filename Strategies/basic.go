package strategies

import (
	helpers "stockApp/Helpers"
)


type Basic struct {
	OpeningPrice float64
	ClosingPrice float64
	Threshold float64
}

func (s *Basic) Validate() bool {
	return (helpers.PercentageDifference(s.OpeningPrice, s.ClosingPrice) >= s.Threshold)  
}