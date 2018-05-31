package openledger

import (
	"github.com/scorum/openledger-go/types"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	testNet = "wss://node.testnet.bitshares.eu"
	mainNet = "wss://bitshares.openledger.info/ws"
)

func TestClient(t *testing.T) {
	t.Run("valid ws url", func(t *testing.T) {
		_, err := NewClient(mainNet)
		require.NoError(t, err)
	})

	t.Run("invalid ws url", func(t *testing.T) {
		_, err := NewClient("wss://invalid")
		require.Error(t, err)
	})
}

func TestClient_Transfer(t *testing.T) {
	t.SkipNow()

	client, err := NewClient(testNet)
	require.NoError(t, err)

	scorum1 := "5KBuq5WmHvgePmB7w3onYsqLM8ESomM2Ae7SigYuuwg8MDHW7NN"
	from := types.MustParseObjectID("1.2.124")
	to := types.MustParseObjectID("1.2.1241")
	amount := types.AssetAmount{
		AssetID: types.MustParseObjectID("1.3.0"),
		Amount:  10000,
	}

	require.NoError(t, client.Transfer(scorum1, from, to, amount))
}

/*
export const limit_order_create = new Serializer("limit_order_create", {
    fee: asset,
    seller: protocol_id_type("account"),
    amount_to_sell: asset,
    min_to_receive: asset,
    expiration: time_point_sec,
    fill_or_kill: bool,
    extensions: set(future_extensions)
});

export const limit_order_cancel = new Serializer("limit_order_cancel", {
    fee: asset,
    fee_paying_account: protocol_id_type("account"),
    order: protocol_id_type("limit_order"),
    extensions: set(future_extensions)
});
*/

// python tests
// https://github.com/bitshares/python-bitshares/blob/9250544ca8eadf66de31c7f38fc37294c11f9548/tests/test_transactions.py
