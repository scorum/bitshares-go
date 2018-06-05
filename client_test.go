package openledger

import (
	"encoding/json"
	"github.com/scorum/openledger-go/sign"
	"github.com/scorum/openledger-go/types"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
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

	scorum1 := "5JiTY3m9u1iPfoKsZdn18pnf26XvX2WnXFJckSiSaiUniNVzxLn"
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

	require.NoError(t, client.Transfer(scorum1, from, to, amount, fee))
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

func TestSign(t *testing.T) {
	t.Skip()
	//test_message := "4018d7844c78f6a6c41c6a552b898022310fc5dec06da467ee7905a8dad512c810b9a4cb8255c93b155b0100b4280000000000000081bc3b95b20140420f000000000000000000"

	op := types.TransferOperation{}
	op.To = types.MustParseObjectID("1.2.22805")
	op.From = types.MustParseObjectID("1.2.974337")
	op.Fee = types.AssetAmount{
		Amount:  10420,
		AssetID: types.MustParseObjectID("1.3.0"),
	}
	op.Amount = types.AssetAmount{
		Amount:  1000000,
		AssetID: types.MustParseObjectID("1.3.0"),
	}
	op.Extensions = []json.RawMessage{}

	tr := types.Transaction{}
	tr.RefBlockPrefix = 1434635172
	tr.RefBlockNum = 47376
	ti := time.Unix(1528118217, 0)
	tr.Expiration = types.Time{&ti}
	tr.Operations = append(tr.Operations, &op)

	st := sign.NewSignedTransaction(&tr)
	_, err := st.Digest("4018d7844c78f6a6c41c6a552b898022310fc5dec06da467ee7905a8dad512c8")
	if err != nil {
		t.Fatal(err)
	}

	//TODO: Actually here we need to copare m
}
