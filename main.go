package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	domain "stockApp/Domain"
	helpers "stockApp/Helpers"
	orders "stockApp/Orders"
	strategies "stockApp/Strategies"
	"time"

	"github.com/coder/websocket"
	"github.com/joho/godotenv"
)

// func main() {
// 	/// sell
// 	log.Print(helpers.PercentageDifference(144, 164.8))
// 	log.Print(helpers.PercentageDifference(144, 164.8) > 5)

// 	// buy
// 	log.Print(helpers.PercentageDifference(174, 154.8))
// 	log.Print(helpers.PercentageDifference(174, 154.8) < -5)
// }

func main() {
	//url := "wss://stream.data.alpaca.markets/v2/test"
	url := "wss://stream.data.alpaca.markets/v2/iex"
	envFile, _ := godotenv.Read(".env")

	keyID, ok := envFile["ALPACA_API_KEY"]
	if !ok {
		log.Fatal("ALPACA_API_KEY not found in .env file")
	}

	secret, ok := envFile["ALPACA_API_SECRET"]
	if !ok {
		log.Fatal("ALPACA_API_SECRET not found in .env file")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	headers := http.Header{
		"APCA-API-KEY-ID":    []string{keyID},
		"APCA-API-SECRET-KEY": []string{secret},
	}
	opts := &websocket.DialOptions{
		HTTPHeader: headers,
	}

	conn, _, err := websocket.Dial(ctx, url, opts)

	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "");

	// we are only interested at trades from the bars aggregate 
	// it will monitor all trades indexes *, can be set to specific index too "AAPL"
	subscribeMsg := `{"action":"subscribe","bars":["*"]}` 

	// --- example of different ways to monitor the feed //
	//subscribeMsg := `{"action":"subscribe","trades":["FAKEPACA"],"quotes":["FAKEPACA"],"bars":["*"]}`
	//subscribeMsg := `{"action":"subscribe","trades":["AAPL","TSLA"],"quotes":["AMD","CLDR"],"bars":["*"]}`

	err = conn.Write(ctx, websocket.MessageBinary, []byte(subscribeMsg) )
	if err != nil {
		log.Fatal("write:", err)
	}

	// TODO, this will set to be since the app is running 
	// find a way to update it every x minutes
	var initialTradeData []domain.Trade
	messageCounter := 0
	for {
		_, d, err := conn.Read(ctx)
		if err != nil {
			log.Println("read:", err)
			break
		}

		messageCounter++

		// IGNORE THE FIRST 3 message as they are for the connection status
		if messageCounter > 3 {
			var currentData []domain.Trade
			errJson := json.Unmarshal(d, &currentData)
			if errJson != nil {
				log.Println(err)
				break
			}

			if len(initialTradeData) == 0 {
				log.Printf("initial data: %v", len(initialTradeData))
				initialTradeData = currentData
			}

			log.Println("current data: ", len(currentData))
			newTradesCached := 0
			for _, trade := range currentData {
				initialTrade := helpers.FilterTradeByIndex(&initialTradeData, trade.Index)
			
				if initialTrade == nil {
					initialTradeData = append(initialTradeData, trade)
					newTradesCached++
					continue
				}

				if trade.Index == initialTrade.Index {
					// TODO: extract as many strategies can run at the same time
					basicBuy := strategies.BasicBuy{
						OpeningPrice: initialTrade.Open,
						ClosingPrice: trade.Close,
						Threshold: -1.5, // %
					}
					
					// Create by order for each index that go below -0.2 %
					order := orders.AlpacaBasic{
						Symbol: trade.Index,
						Price: trade.Close,
						Quantity: 10,
					}
				
					if(basicBuy.Validate()) {
						order.Buy()
					} 
					/// SELL ALL WITH PROFIOT OF 2.5%
					basicSell := strategies.BasicSell{
						OpeningPrice: initialTrade.Open,
						ClosingPrice: trade.Close,
						Threshold: 2.5, // %
					}

					if (basicSell.Validate()) {
						order.Sell()
					}
				}
			}
			
			if (newTradesCached > 0) {
				log.Printf("TotalCached %v  New Added: %v", len(initialTradeData), newTradesCached)
				newTradesCached = 0
			}

		} else {
			log.Printf("Received message: %s\n", d)
		}
		
	}
}