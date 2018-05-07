package history

import "github.com/scorum/openledger-go/types"

type Bucket struct {
	ID          string `json:"id"`
	Key         Key    `json:"key"`
	HighBase    uint32 `json:"high_base"`
	HighQuote   uint32 `json:"high_quote"`
	LowBase     uint32 `json:"low_base"`
	LowQuote    uint32 `json:"low_quote"`
	OpenBase    uint32 `json:"open_base"`
	OpenQuote   uint32 `json:"open_quote"`
	CloseBase   uint32 `json:"close_base"`
	CloseQuote  uint32 `json:"close_quote"`
	BaseVolume  uint32 `json:"base_volume"`
	QuoteVolume uint32 `json:"quote_volume"`
}

type Key struct {
	Base    string     `json:"base"`
	Quote   string     `json:"quote"`
	Seconds uint32     `json:"seconds"`
	Open    types.Time `json:"open"`
}
