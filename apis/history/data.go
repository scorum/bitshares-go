package history

import (
	"encoding/json"
	"github.com/scorum/bitshares-go/types"
)

type OrderHistory struct {
	ID        string         `json:"id"`
	Key       OrderKey       `json:"key"`
	Time      types.Time     `json:"time"`
	Operation OrderOperation `json:"op"`
}

type OrderKey struct {
	Base     types.ObjectID `json:"base"`
	Quote    types.ObjectID `json:"quote"`
	Sequence int32          `json:"sequence"`
}

type OrderOperation struct {
	Fee       types.AssetAmount `json:"fee"`
	Pays      types.AssetAmount `json:"pays"`
	Receives  types.AssetAmount `json:"receives"`
	FillPrice types.Price       `json:"fill_price"`
	OrderID   types.ObjectID    `json:"order_id"`
	AccountID types.ObjectID    `json:"account_id"`
	IsMaker   bool              `json:"is_maker"`
}

type OperationHistory struct {
	ID                       string            `json:"id"`
	BlockNumber              uint32            `json:"block_num"`
	TransactionsInBlock      uint32            `json:"trx_in_block"`
	OperationsInTransactions uint32            `json:"op_in_trx"`
	VirtualOperations        uint32            `json:"virtual_op"`
	Result                   []json.RawMessage `json:"result"`
	Operations               []json.RawMessage `json:"op"`
}
