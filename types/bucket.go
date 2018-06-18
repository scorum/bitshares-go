package types

type Bucket struct {
	ID          string    `json:"id"`
	Key         BucketKey `json:"key"`
	HighBase    Suint32   `json:"high_base"`
	HighQuote   Suint32   `json:"high_quote"`
	LowBase     Suint32   `json:"low_base"`
	LowQuote    Suint32   `json:"low_quote"`
	OpenBase    Suint32   `json:"open_base"`
	OpenQuote   Suint32   `json:"open_quote"`
	CloseBase   Suint32   `json:"close_base"`
	CloseQuote  Suint32   `json:"close_quote"`
	BaseVolume  Suint32   `json:"base_volume"`
	QuoteVolume Suint32   `json:"quote_volume"`
}

type BucketKey struct {
	Base    string  `json:"base"`
	Quote   string  `json:"quote"`
	Seconds Suint32 `json:"seconds"`
	Open    Time    `json:"open"`
}
