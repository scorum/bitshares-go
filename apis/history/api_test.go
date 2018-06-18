package history

import (
	"github.com/scorum/bitshares-go/apis/database"
	"github.com/scorum/bitshares-go/apis/login"
	"github.com/scorum/bitshares-go/transport/websocket"
	"github.com/scorum/bitshares-go/types"
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

	// [60,900,1800,3600,14400,86400] in seconds
	//require.Len(t, buckets, 7)

	require.NotEmpty(t, buckets)
}

func TestGetFillOrderHistory(t *testing.T) {
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

	orders, err := historyAPI.GetFillOrderHistory(openSCR, USD, 5)
	require.True(t, len(orders) > 0)
}

func TestAccountHistory(t *testing.T) {
	transport, err := websocket.NewTransport(url)
	require.NoError(t, err)
	defer transport.Close()

	// request access to the history api
	historyAPIID, err := login.NewAPI(transport).History()
	require.NoError(t, err)

	historyAPI := NewAPI(historyAPIID, transport)

	user, _ := types.ParseObjectID("1.2.900546")
	start, _ := types.ParseObjectID("1.11.0")
	stop, _ := types.ParseObjectID("1.11.0")

	history, err := historyAPI.GetAccountHistory(user, stop, 100, start)
	require.NoError(t, err)
	require.NotEmpty(t, history)
}
