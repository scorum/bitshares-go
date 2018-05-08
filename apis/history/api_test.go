package history

import (
	"github.com/scorum/openledger-go/apis/database"
	"github.com/scorum/openledger-go/apis/login"
	"github.com/scorum/openledger-go/transport/websocket"
	"github.com/scorum/openledger-go/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const url = "wss://bitshares.openledger.info/ws"

func TestGetMarketHistory(t *testing.T) {
	transport, err := websocket.NewTransport(url)
	require.NoError(t, err)
	defer transport.Close()

	// request access to the database api
	databaseAPIID, err := login.NewAPI(transport).Database()
	require.NoError(t, err)

	// request access to the history api
	historyAPIID, err := login.NewAPI(transport).History()
	require.NoError(t, err)

	historyAPI := NewAPI(historyAPIID, transport)
	databaseAPI := database.NewAPI(databaseAPIID, transport)

	// lookup symbols ids
	symbols, err := databaseAPI.LookupAssetSymbols("OPEN.SCR", "USD")
	require.NoError(t, err)

	openSCR := symbols[0].ID
	USD := symbols[1].ID

	day := time.Hour * 24
	start := types.NewTime(time.Now().Add(-5 * day))
	end := types.NewTime(time.Now())
	buckets, err := historyAPI.GetMarketHistory(openSCR, USD, 60, start, end)

	// assert
	require.NoError(t, err)
	require.NotNil(t, buckets)
}

func TestGetMarketHistoryBuckets(t *testing.T) {
	transport, err := websocket.NewTransport(url)
	require.NoError(t, err)
	defer transport.Close()

	// request access to the history api
	historyAPIID, err := login.NewAPI(transport).History()
	require.NoError(t, err)

	historyAPI := NewAPI(historyAPIID, transport)

	buckets, err := historyAPI.GetMarketHistoryBuckets()
	require.NoError(t, err)

	// [15,60,300,3600,86400] in seconds
	require.Len(t, buckets, 5)
}
