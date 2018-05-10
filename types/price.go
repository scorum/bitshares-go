package types

import (
	"encoding/json"
	"strconv"
)

type Price struct {
	Base  AssetAmount `json:"base"`
	Quote AssetAmount `json:"quote"`
}

type AssetAmount struct {
	Amount  uint64   `json:"amount"`
	AssetID ObjectID `json:"asset_id"`
}

// RPC client might return asset amount as uint64 or string,
// therefore a custom unmarshaller is used
func (ops *AssetAmount) UnmarshalJSON(b []byte) (err error) {
	stringCase := struct {
		Amount  string   `json:"amount"`
		AssetID ObjectID `json:"asset_id"`
	}{}

	uint64Case := struct {
		Amount  uint64   `json:"amount"`
		AssetID ObjectID `json:"asset_id"`
	}{}

	if err = json.Unmarshal(b, &uint64Case); err == nil {
		ops.AssetID = uint64Case.AssetID
		ops.Amount = uint64Case.Amount
		return nil
	}

	// failed on uint64, try string
	if err = json.Unmarshal(b, &stringCase); err == nil {
		ops.AssetID = stringCase.AssetID
		ops.Amount, err = strconv.ParseUint(stringCase.Amount, 10, 64)
		return err
	}

	return err
}
