package orders

import (
	"log"
	helpers "stockApp/Helpers"
	"sync"
)

type AlpacaBasic struct {
	Symbol string
	Price float64
	Quantity int
}

func (a *AlpacaBasic) Buy() bool {
	//todo implement the logic for alpaca buy
	var wg sync.WaitGroup
  wg.Add(2)
  go func() {
    helpers.PlaySound(helpers.MONEY_OUT)
    wg.Done()
  }()
  go func() {
    log.Printf("%s 💸💸💸💸", a.Symbol)
    wg.Done()
  }()
  wg.Wait()
	return false
}


func (a *AlpacaBasic) Sell() bool {
	//todo implement the logic for alpaca sell
	helpers.PlaySound(helpers.MONEY_IN)
	helpers.PrintCashIn()
	log.Printf("%s SOLD", a.Symbol)
	return false
}