package types

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAssetAmount_UnmarshalJSON(t *testing.T) {
	t.Run("amount uint64", func(t *testing.T) {
		data := `{
              "amount": 11211,
              "asset_id": "1.3.0"
            }`
		am := AssetAmount{}
		require.NoError(t, json.Unmarshal([]byte(data), &am))

		require.Equal(t, uint64(11211), am.Amount)
		require.Equal(t, ObjectID{Space: 1, Type: 3}, am.AssetID)
	})

	t.Run("amount string", func(t *testing.T) {
		data := `{
              "amount": "14450706212",
              "asset_id": "1.3.3232"
            }`
		am := AssetAmount{}
		require.NoError(t, json.Unmarshal([]byte(data), &am))

		require.Equal(t, uint64(14450706212), am.Amount)
		require.Equal(t, ObjectID{Space: 1, Type: 3, ID: 3232}, am.AssetID)
	})
}
