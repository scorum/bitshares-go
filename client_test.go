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
	client, err := NewClient(testNet)
	require.NoError(t, err)

	scorum1 := "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"
	from := types.MustParseObjectID("1.2.124")
	to := types.MustParseObjectID("1.2.1241")
	amount := types.AssetAmount{
		AssetID: types.MustParseObjectID("1.3.0"),
		Amount:  10000,
	}

	require.NoError(t, client.Transfer(scorum1, from, to, amount))
}
