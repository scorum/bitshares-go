package database

import (
	"encoding/json"
	"github.com/scorum/bitshares-go/caller"
	"github.com/scorum/bitshares-go/types"
)

type API struct {
	caller caller.Caller
	id     caller.APIID
}

func NewAPI(id caller.APIID, caller caller.Caller) *API {
	return &API{id: id, caller: caller}
}

func (api *API) call(method string, args []interface{}, reply interface{}) error {
	return api.caller.Call(api.id, method, args, reply)
}

func (api *API) setCallback(method string, callback func(raw json.RawMessage)) error {
	return api.caller.SetCallback(api.id, method, callback)
}

func (api *API) GetChainID() (*string, error) {
	var resp string
	err := api.call("get_chain_id", caller.EmptyParams, &resp)
	return &resp, err
}

// GetConfig retrieves compile-time constants
func (api *API) GetConfig() (*Config, error) {
	var config Config
	err := api.call("get_config", caller.EmptyParams, &config)
	return &config, err
}

// GetTransaction used to fetch an individual transaction
func (api *API) GetTransaction(blockNum uint32, trxInBlock uint32) (*types.Transaction, error) {
	var resp types.Transaction
	err := api.call("get_transaction", []interface{}{blockNum, trxInBlock}, &resp)
	return &resp, err
}

// GetRecentTransactionByID
// If the transaction has not expired, this method will return the transaction for the given ID or
// it will return NULL if it is not known. Just because it is not known does not mean
// it wasn’t included in the blockchain.
func (api *API) GetRecentTransactionByID(transactionID uint32) (*types.Transaction, error) {
	var resp types.Transaction
	err := api.call("get_recent_transaction_by_id", []interface{}{transactionID}, &resp)
	return &resp, err
}

// GetDynamicGlobalProperties retrieves the current global_property_object
func (api *API) GetDynamicGlobalProperties() (*DynamicGlobalProperties, error) {
	var resp DynamicGlobalProperties
	err := api.call("get_dynamic_global_properties", caller.EmptyParams, &resp)
	return &resp, err
}

// LookupAssetSymbols get assets corresponding to the provided symbols or IDs
func (api *API) LookupAssetSymbols(symbols ...string) ([]*Asset, error) {
	var resp []*Asset
	err := api.call("lookup_asset_symbols", []interface{}{symbols}, &resp)
	return resp, err
}

// GetLimitOrders returns limit orders in a given market.
// There are both sell and buy orders.
// For the sell orders LimitOrder.SellPrice.Base = the given base
// For the buy orders LimitOrder.SellPrice.Base = the given quote
func (api *API) GetLimitOrders(base, quote types.ObjectID, limit uint32) ([]*LimitOrder, error) {
	var resp []*LimitOrder
	err := api.call("get_limit_orders", []interface{}{base.String(), quote.String(), limit}, &resp)
	return resp, err
}

// GetBlockHeader returns block header by the given block number
func (api *API) GetBlockHeader(blockNum uint32) (*BlockHeader, error) {
	var resp BlockHeader
	err := api.call("get_block_header", []interface{}{blockNum}, &resp)
	return &resp, err
}

// GetBlock return a block by the given block number
func (api *API) GetBlock(blockNum uint32) (*Block, error) {
	var resp Block
	err := api.call("get_block", []interface{}{blockNum}, &resp)
	return &resp, err
}

// GetBlock return a block by the given block number
func (api *API) GetObjects(assets ...types.ObjectID) ([]json.RawMessage, error) {
	var resp []json.RawMessage
	err := api.call("get_objects", []interface{}{objectsToParams(assets)}, &resp)
	return resp, err
}

// GetTicker returns the ticker for the market assetA:assetB (past 24 hours)
func (api *API) GetTicker(base, quote types.ObjectID) (*MarketTicker, error) {
	var resp MarketTicker
	err := api.call("get_ticker", []interface{}{base.String(), quote.String()}, &resp)
	return &resp, err
}

// GetAccountBalances
// Get an account’s balances in various assets.
func (api *API) GetAccountBalances(accountID types.ObjectID, assets ...types.ObjectID) ([]*types.AssetAmount, error) {
	var resp []*types.AssetAmount
	err := api.call("get_account_balances", []interface{}{accountID.String(), objectsToParams(assets)}, &resp)
	return resp, err
}

func objectsToParams(objs []types.ObjectID) []string {
	objsStr := make([]string, len(objs))
	for i, asset := range objs {
		objsStr[i] = asset.String()
	}
	return objsStr
}

// Semantically equivalent to get_account_balances, but takes a name instead of an ID.
func (api *API) GetNamedAccountBalances(account string, assets ...types.ObjectID) ([]*types.AssetAmount, error) {
	var resp []*types.AssetAmount
	err := api.call("get_named_account_balances", []interface{}{account, objectsToParams(assets)}, &resp)
	return resp, err
}

// LookupAccounts gets names and IDs for registered accounts
// lower_bound_name: Lower bound of the first name to return
// limit: Maximum number of results to return must not exceed 1000
func (api *API) LookupAccounts(lowerBoundName string, limit uint16) (AccountsMap, error) {
	var resp AccountsMap
	err := api.call("lookup_accounts", []interface{}{lowerBoundName, limit}, &resp)
	return resp, err
}

// SetBlockAppliedCallback registers a global subscription callback
func (api *API) SetBlockAppliedCallback(notice func(blockID string, err error)) (err error) {
	err = api.setCallback("set_block_applied_callback", func(raw json.RawMessage) {
		var header []string
		if err := json.Unmarshal(raw, &header); err != nil {
			notice("", err)
		}
		for _, b := range header {
			notice(b, nil)
		}
	})
	return
}

// CancelAllSubscriptions cancel all subscriptions
func (api *API) CancelAllSubscriptions() error {
	return api.call("cancel_all_subscriptions", caller.EmptyParams, nil)
}

// GetRequiredFee fetchs fee for operations
func (api *API) GetRequiredFee(ops []types.Operation, assetID string) ([]types.AssetAmount, error) {
	var resp []types.AssetAmount

	opsJSON := []interface{}{}
	for _, o := range ops {
		_, err := json.Marshal(o)
		if err != nil {
			return []types.AssetAmount{}, err
		}

		opArr := []interface{}{o.Type(), o}

		opsJSON = append(opsJSON, opArr)
	}

	err := api.call("get_required_fees", []interface{}{opsJSON, assetID}, &resp)
	return resp, err
}
