package main

import (
	"context"
	"log"
	"net/http"
	orders "stockApp/Orders"
	strategies "stockApp/Strategies"
	"time"

	"github.com/coder/websocket"
	"github.com/joho/godotenv"
)

func main() {
	basicStrategy := strategies.Basic{
		OpeningPrice: 100,
		ClosingPrice: 80,
		Threshold: 20,
	}

	order := orders.AlpacaBasic{
		Symbol: "AAPL",
		Price: 100,
		Quantity: 10,
	}

	validStrategy := basicStrategy.Validate()

	if(validStrategy) {
		order.Buy()
	}else {
		order.Sell()
	}
	return
	url := "wss://stream.data.alpaca.markets/v2/test"
	//url := "wss://stream.data.alpaca.markets/v2/iex"
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

	subscribeMsg := `{"action":"subscribe","bars":["*"]}`
	//subscribeMsg := `{"action":"subscribe","trades":["FAKEPACA"],"quotes":["FAKEPACA"],"bars":["*"]}`
	//subscribeMsg := `{"action":"subscribe","trades":["AAPL","TSLA"],"quotes":["AMD","CLDR"],"bars":["*"]}`

	err = conn.Write(ctx, websocket.MessageBinary, []byte(subscribeMsg) )
	if err != nil {
		log.Fatal("write:", err)
	}

	for {
		_,b, err := conn.Read(ctx)
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("Received message: %s\n", b)
	}
}