package openledger

import (
	"github.com/scorum/openledger-go/apis/database"
	"github.com/scorum/openledger-go/apis/history"
	"github.com/scorum/openledger-go/apis/login"
	"github.com/scorum/openledger-go/apis/networkbroadcast"
	"github.com/scorum/openledger-go/caller"
	"github.com/scorum/openledger-go/transport/websocket"
)

type Client struct {
	cc caller.CallCloser

	// Database represents database_api
	Database *database.API

	// NetworkBroadcast represents network_broadcast_api
	NetworkBroadcast *networkbroadcast.API

	// History represents history_api
	History *history.API

	// Login represents login_api
	Login *login.API
}

// NewClient creates a new RPC client
func NewClient(url string) (*Client, error) {
	// transport
	transport, err := websocket.NewTransport(url)
	if err != nil {
		return nil, err
	}

	client := &Client{cc: transport}

	// login
	loginAPI := login.NewAPI(transport)
	client.Login = loginAPI

	// database
	databaseAPIID, err := loginAPI.Database()
	if err != nil {
		return nil, err
	}
	client.Database = database.NewAPI(databaseAPIID, client.cc)

	// history
	historyAPIID, err := loginAPI.History()
	if err != nil {
		return nil, err
	}
	client.History = history.NewAPI(historyAPIID, client.cc)

	// network broadcast
	networkBroadcastAPIID, err := loginAPI.NetworkBroadcast()
	if err != nil {
		return nil, err
	}
	client.NetworkBroadcast = networkbroadcast.NewAPI(networkBroadcastAPIID, client.cc)

	return client, nil
}

// Close should be used to close the client when no longer needed.
// It simply calls Close() on the underlying CallCloser.
func (client *Client) Close() error {
	return client.cc.Close()
}
