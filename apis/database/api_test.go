package database

import (
	"github.com/scorum/openledger-go/apis/login"
	"github.com/scorum/openledger-go/transport/websocket"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const url = "wss://bitshares.openledger.info/ws"

func getAPI(t *testing.T) *API {
	transport, err := websocket.NewTransport(url)
	require.NoError(t, err)

	// request access to the database api
	databaseAPIID, err := login.NewAPI(transport).Database()
	require.NoError(t, err)

	return NewAPI(*databaseAPIID, transport)
}

func TestGetChainID(t *testing.T) {
	databaseAPI := getAPI(t)
	chainID, err := databaseAPI.GetChainID()

	// assert
	require.NoError(t, err)
	require.NotNil(t, chainID)
}

func TestLookupAssetSymbol(t *testing.T) {
	databaseAPI := getAPI(t)
	symbols, err := databaseAPI.LookupAssetSymbols("OPEN.BTC", "USD")
	require.NoError(t, err)

	require.Len(t, symbols, 2)
	require.Equal(t, "OPEN.BTC", symbols[0].Symbol)
	require.Equal(t, "USD", symbols[1].Symbol)

	require.Equal(t, "1.3.861", symbols[0].ID.String())
	require.Equal(t, "1.3.121", symbols[1].ID.String())
}

func TestGetBlockHeader(t *testing.T) {
	databaseAPI := getAPI(t)
	header, err := databaseAPI.GetBlockHeader(25)

	require.NoError(t, err)
	require.NotEmpty(t, header.Previous)
	require.NotEmpty(t, header.Witness)
}

func TestGetTicker(t *testing.T) {
	databaseAPI := getAPI(t)
	symbols, err := databaseAPI.LookupAssetSymbols("OPEN.SCR", "USD")
	require.NoError(t, err)

	openSCR := symbols[0].ID
	USD := symbols[1].ID

	ticker, err := databaseAPI.GetTicker(openSCR, USD)
	require.NoError(t, err)

	require.Equal(t, ticker.Base, openSCR)
	require.Equal(t, ticker.Quote, USD)
}

func TestSetBlockAppliedCallback(t *testing.T) {
	databaseAPI := getAPI(t)

	var called bool
	err := databaseAPI.SetBlockAppliedCallback(func(blockID string, err error) {
		t.Log("block:", blockID)
		require.NoError(t, err)
		called = true
	})
	require.NoError(t, err)
	time.Sleep(10 * time.Second)
	require.True(t, called)

	require.NoError(t, databaseAPI.CancelAllSubscriptions())
}
