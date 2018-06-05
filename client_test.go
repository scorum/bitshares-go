package openledger

import (
	"github.com/scorum/openledger-go/types"
	"github.com/stretchr/testify/require"
	"log"
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
	client, err := NewClient(testNet)
	require.Nil(t, err)

	cali4888arr, err := client.Database.LookupAccounts("cali4889", 2)
	require.Nil(t, err)

	log.Println(cali4888arr["cali4889"])

	cali4889ID := cali4888arr["cali4889"]
	cali4890ID := cali4888arr["cali4890"]

	assets, err := client.Database.LookupAssetSymbols("TEST")
	require.Nil(t, err)

	cali4889IDActiveKey := "5JiTY3m9u1iPfoKsZdn18pnf26XvX2WnXFJckSiSaiUniNVzxLn"
	from := cali4889ID
	to := cali4890ID
	amount := types.AssetAmount{
		AssetID: assets[0].ID,
		Amount:  1000,
	}
	fee := types.AssetAmount{
		AssetID: assets[0].ID,
		Amount:  0,
	}

	require.NoError(t, client.Transfer(cali4889IDActiveKey, from, to, amount, fee))
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
