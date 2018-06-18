package bitshares

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/scorum/bitshares-go/apis/database"
	"github.com/scorum/bitshares-go/apis/history"
	"github.com/scorum/bitshares-go/apis/login"
	"github.com/scorum/bitshares-go/apis/networkbroadcast"
	"github.com/scorum/bitshares-go/caller"
	"github.com/scorum/bitshares-go/sign"
	"github.com/scorum/bitshares-go/transport/websocket"
	"github.com/scorum/bitshares-go/types"
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

	stx, err := client.sign([]string{key}, op)
	if err != nil {
		return err
	}
	return client.broadcast(stx)
}

func (client *Client) LimitOrderCreate(key string, seller types.ObjectID, fee, amToSell, minToRecive types.AssetAmount, expiration time.Duration, fillOrKill bool) (string, error) {
	props, err := client.Database.GetDynamicGlobalProperties()
	if err != nil {
		return "", errors.Wrap(err, "failed to get dynamic global properties")
	}

	op := &types.LimitOrderCreateOperation{
		Fee:          fee,
		Seller:       seller,
		AmountToSell: amToSell,
		MinToReceive: minToRecive,
		Expiration:   types.NewTime(props.Time.Add(expiration)),
		FillOrKill:   fillOrKill,
		Extensions:   []json.RawMessage{},
	}

	fees, err := client.Database.GetRequiredFee([]types.Operation{op}, fee.AssetID.String())
	if err != nil {
		log.Println(err)
		return "", errors.Wrap(err, "can't get fees")
	}
	op.Fee.Amount = fees[0].Amount

	stx, err := client.sign([]string{key}, op)
	if err != nil {
		return "", err
	}
	result, err := client.broadcastSync(stx)
	if err != nil {
		return "", err
	}

	res := result.Trx["operation_results"]
	ops, ok := res.([]interface{})
	if !ok {
		return "", errors.New("invalid result format")
	}
	create_op, ok := ops[0].([]interface{})
	if !ok {
		return "", errors.New("invalid result format")
	}
	id, ok := create_op[1].(string)
	if !ok {
		return "", errors.New("invalid result format")
	}

	return id, err
}

func (client *Client) LimitOrderCancel(key string, feePayingAccount, order types.ObjectID, fee types.AssetAmount) error {
	op := &types.LimitOrderCancelOperation{
		Fee:              fee,
		FeePayingAccount: feePayingAccount,
		Order:            order,
		Extensions:       []json.RawMessage{},
	}

	fees, err := client.Database.GetRequiredFee([]types.Operation{op}, fee.AssetID.String())
	if err != nil {
		log.Println(err)
		return errors.Wrap(err, "can't get fees")
	}
	op.Fee.Amount = fees[0].Amount

	stx, err := client.sign([]string{key}, op)
	if err != nil {
		return err
	}
	return client.broadcast(stx)
}

func (client *Client) sign(wifs []string, operations ...types.Operation) (*sign.SignedTransaction, error) {
	props, err := client.Database.GetDynamicGlobalProperties()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get dynamic global properties")
	}

	block, err := client.Database.GetBlock(props.LastIrreversibleBlockNum)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get block")
	}

	refBlockPrefix, err := sign.RefBlockPrefix(block.Previous)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign block prefix")
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
		return nil, errors.Wrap(err, "failed to sign the transaction")
	}

	return stx, nil
}

func (client *Client) broadcast(stx *sign.SignedTransaction) error {
	return client.NetworkBroadcast.BroadcastTransaction(stx.Transaction)
}

func (client *Client) broadcastSync(stx *sign.SignedTransaction) (*networkbroadcast.BroadcastResponse, error) {
	return client.NetworkBroadcast.BroadcastTransactionSynchronous(stx.Transaction)
}
