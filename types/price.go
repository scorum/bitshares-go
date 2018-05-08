package types

type Price struct {
	Base  AssetAmount `json:"base"`
	Quote AssetAmount `json:"quote"`
}

type AssetAmount struct {
	Amount  uint64   `json:"amount"`
	AssetID ObjectID `json:"asset_id"`
}
