package history

import (
	"github.com/scorum/openledger-go/caller"
	"github.com/scorum/openledger-go/types"
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

// GetMarketHistory
func (api *API) GetMarketHistory(base, quote types.ObjectID, bucketSeconds uint32, start, end types.Time) ([]*Bucket, error) {
	var resp []*Bucket
	err := api.call("get_market_history", []interface{}{base.String(), quote.String(), bucketSeconds, start, end}, &resp)
	return resp, err
}
