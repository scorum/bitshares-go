package database

import (
	"github.com/scorum/openledger-go/apis/login"
	"github.com/scorum/openledger-go/transport/websocket"
	"github.com/scorum/openledger-go/types"
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

	return NewAPI(databaseAPIID, transport)
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

func TestGetAccountBalances(t *testing.T) {
	databaseAPI := getAPI(t)
	user, _ := types.ParseObjectID("1.2.900546")

	t.Run("empty assets should return all balanced", func(t *testing.T) {
		balances, err := databaseAPI.GetAccountBalances(user)
		require.NoError(t, err)
		require.NotEmpty(t, balances)
	})

	t.Run("request SCR balance", func(t *testing.T) {
		symbols, _ := databaseAPI.LookupAssetSymbols("OPEN.SCR")

		balances, err := databaseAPI.GetAccountBalances(user, symbols[0].ID)
		require.NoError(t, err)
		require.NotEmpty(t, balances)
	})

	t.Run("request USD balance", func(t *testing.T) {
		symbols, _ := databaseAPI.LookupAssetSymbols("USD")

		balances, err := databaseAPI.GetAccountBalances(user, symbols[0].ID)
		require.NoError(t, err)
		require.NotEmpty(t, balances)
		require.Equal(t, uint64(0), balances[0].Amount)
	})
}

func TestGetNamedAccountBalances(t *testing.T) {
	databaseAPI := getAPI(t)

	t.Run("empty assets should return all balanced", func(t *testing.T) {
		balances, err := databaseAPI.GetNamedAccountBalances("megaherz1")
		require.NoError(t, err)
		require.NotEmpty(t, balances)
	})

	t.Run("not existing account", func(t *testing.T) {
		_, err := databaseAPI.GetNamedAccountBalances("nonexists")
		require.Error(t, err)
	})
}

func TestLookupAccounts(t *testing.T) {
	databaseAPI := getAPI(t)

	t.Run("empty lower bound", func(t *testing.T) {
		accounts, err := databaseAPI.LookupAccounts("", 3)
		require.NoError(t, err)
		require.Len(t, accounts, 3)
	})

	t.Run("limit exceeded", func(t *testing.T) {
		_, err := databaseAPI.LookupAccounts("", 1001)
		require.Error(t, err)
	})
}

func TestGetLimitOrders(t *testing.T) {
	databaseAPI := getAPI(t)
	symbols, err := databaseAPI.LookupAssetSymbols("OPEN.BTC", "USD")
	require.NoError(t, err)

	orders, err := databaseAPI.GetLimitOrders(symbols[0].ID, symbols[1].ID, 100)
	require.NoError(t, err)
	require.NotEmpty(t, orders)
}

func TestGetBlockHeader(t *testing.T) {
	databaseAPI := getAPI(t)
	header, err := databaseAPI.GetBlockHeader(25)

	require.NoError(t, err)
	require.NotEmpty(t, header.Previous)
	require.NotEmpty(t, header.Witness)
}

func TestGetBlock(t *testing.T) {
	databaseAPI := getAPI(t)
	block, err := databaseAPI.GetBlock(26851093)
	require.NoError(t, err)
	require.NotEmpty(t, block.Previous)
	require.NotEmpty(t, block.Witness)
}

func TestGetTransaction(t *testing.T) {
	databaseAPI := getAPI(t)

	t.Run("first transaction in a block", func(t *testing.T) {
		trx, err := databaseAPI.GetTransaction(26851092, 1)
		require.NoError(t, err)
		require.NotNil(t, trx)
		require.Len(t, trx.Operations, 1)
	})

	t.Run("trx num exceeded", func(t *testing.T) {
		_, err := databaseAPI.GetTransaction(26851092, 100)
		require.Error(t, err) //Assert Exception: opt_block->transactions.size() > trx_num
	})
}

func TestGetRecentTransactionByID(t *testing.T) {
	databaseAPI := getAPI(t)
	trx, err := databaseAPI.GetRecentTransactionByID(3)
	require.NoError(t, err)
	require.NotNil(t, trx)
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

func TestDynamicGlobalProperties(t *testing.T) {
	databaseAPI := getAPI(t)
	props, err := databaseAPI.GetDynamicGlobalProperties()
	require.NoError(t, err)
	require.True(t, props.HeadBlockNumber > 0)
	require.True(t, props.LastIrreversibleBlockNum > 0)
}

func TestGetConfig(t *testing.T) {
	databaseAPI := getAPI(t)
	config, err := databaseAPI.GetConfig()
	require.NoError(t, err)
	require.Equal(t, "BTS", config.GrapheneAddressPrefix)
	require.Equal(t, "BTS", config.GrapheneSymbol)
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
	time.Sleep(5 * time.Second)
	require.True(t, called)

	require.NoError(t, databaseAPI.CancelAllSubscriptions())
}
