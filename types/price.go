package types

import (
	"encoding/json"
	"github.com/scorum/openledger-go/encoding/transaction"
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

func (aa AssetAmount) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(aa.Amount)
	enc.Encode(aa.AssetID)
	return enc.Err()
}

// RPC client might return asset amount as uint64 or string,
// therefore a custom unmarshaller is used
func (aa *AssetAmount) UnmarshalJSON(b []byte) (err error) {
	stringCase := struct {
		Amount  string   `json:"amount"`
		AssetID ObjectID `json:"asset_id"`
	}{}

	uint64Case := struct {
		Amount  uint64   `json:"amount"`
		AssetID ObjectID `json:"asset_id"`
	}{}

	if err = json.Unmarshal(b, &uint64Case); err == nil {
		aa.AssetID = uint64Case.AssetID
		aa.Amount = uint64Case.Amount
		return nil
	}

	// failed on uint64, try string
	if err = json.Unmarshal(b, &stringCase); err == nil {
		aa.AssetID = stringCase.AssetID
		aa.Amount, err = strconv.ParseUint(stringCase.Amount, 10, 64)
		return err
	}

	return err
}
