package openledger

import (
	"github.com/pkg/errors"
	"github.com/scorum/openledger-go/apis/database"
	"github.com/scorum/openledger-go/apis/history"
	"github.com/scorum/openledger-go/apis/login"
	"github.com/scorum/openledger-go/apis/networkbroadcast"
	"github.com/scorum/openledger-go/caller"
	"github.com/scorum/openledger-go/sign"
	"github.com/scorum/openledger-go/transport/websocket"
	"github.com/scorum/openledger-go/types"
	"log"
	"time"
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

	chainID string
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

	// chain ID
	chainID, err := client.Database.GetChainID()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get chain ID")
	}
	client.chainID = *chainID

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

// Transfer a certain amount of the given asset
func (client *Client) Transfer(key string, from, to types.ObjectID, amount, fee types.AssetAmount) error {
	op := types.NewTransferOperation(from, to, amount, fee)

	fees, err := client.Database.GetRequiredFee([]types.Operation{op}, fee.AssetID.String())
	if err != nil {
		log.Println(err)
		return errors.Wrap(err, "can't get fees")
	}
	op.Fee.Amount = fees[0].Amount

	return client.broadcast([]string{key}, op)
}

// Sign the given operations with the wifs and broadcast them as one transaction
func (client *Client) broadcast(wifs []string, operations ...types.Operation) error {
	props, err := client.Database.GetDynamicGlobalProperties()
	if err != nil {
		return errors.Wrap(err, "failed to get dynamic global properties")
	}

	block, err := client.Database.GetBlock(props.LastIrreversibleBlockNum)
	if err != nil {
		return errors.Wrap(err, "failed to get block")
	}

	refBlockPrefix, err := sign.RefBlockPrefix(block.Previous)
	if err != nil {
		return errors.Wrap(err, "failed to sign block prefix")
	}

	expiration := props.Time.Add(10 * time.Minute)
	stx := sign.NewSignedTransaction(&types.Transaction{
		RefBlockNum:    sign.RefBlockNum(props.LastIrreversibleBlockNum - 1&0xffff),
		RefBlockPrefix: refBlockPrefix,
		Expiration:     types.Time{Time: &expiration},
	})

	for _, op := range operations {
		stx.PushOperation(op)
	}

	if err = stx.Sign(wifs, client.chainID); err != nil {
		return errors.Wrap(err, "failed to sign the transaction")
	}

	return client.NetworkBroadcast.BroadcastTransaction(stx.Transaction)
}
