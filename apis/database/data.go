package database

import (
	"encoding/json"
	"github.com/scorum/openledger-go/types"
)

type Asset struct {
	ID                 types.ObjectID `json:"id"`
	Symbol             string         `json:"symbol"`
	Precision          uint8          `json:"precision"`
	Issuer             string         `json:"issuer"`
	DynamicAssetDataID string         `json:"dynamic_asset_data_id"`
}

type BlockHeader struct {
	TransactionMerkleRoot string            `json:"transaction_merkle_root"`
	Previous              string            `json:"previous"`
	Timestamp             types.Time        `json:"timestamp"`
	Witness               string            `json:"witness"`
	Extensions            []json.RawMessage `json:"extensions"`
}

type MarketTicker struct {
	Time          types.Time     `json:"time"`
	Base          types.ObjectID `json:"base"`
	Quote         types.ObjectID `json:"quote"`
	Latest        string         `json:"latest"`
	LowestAsk     string         `json:"lowest_ask"`
	HighestBid    string         `json:"highest_bid"`
	PercentChange string         `json:"percent_change"`
	BaseVolume    string         `json:"base_volume"`
	QuoteVolume   string         `json:"quote_volume"`
}

type LimitOrder struct {
	ID          types.ObjectID `json:"id"`
	Expiration  types.Time     `json:"expiration"`
	Seller      types.ObjectID `json:"seller"`
	ForSale     uint64         `json:"for_sale"`
	DeferredFee uint64         `json:"deferred_fee"`
	SellPrice   types.Price    `json:"sell_price"`
}
