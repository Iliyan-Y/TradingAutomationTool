package domain

type Trade struct {
	Index  string  `json:"S"`
	Open   float64 `json:"o"`
	High   float64 `json:"h"`
	Low    float64 `json:"l"`
	Close  float64 `json:"c"`
	Volume int     `json:"v"`
	Timestamp    string  `json:"t"`
	MessageType string  `json:"T"`
}
