package login

import "github.com/scorum/bitshares-go/caller"

const APIID = 1

type API struct {
	caller caller.Caller
}

func NewAPI(caller caller.Caller) *API {
	return &API{caller}
}

func (api *API) call(method string, args []interface{}, reply interface{}) error {
	return api.caller.Call(caller.APIID(APIID), method, args, reply)
}

func (api *API) GetApiByName(name string) (*uint8, error) {
	var id uint8
	err := api.call("get_api_by_name", []interface{}{name}, &id)
	return &id, err
}

func (api *API) Login(username, password string) (bool, error) {
	var resp bool
	err := api.call("login", []interface{}{username, password}, &resp)
	return resp, err
}

func (api *API) Database() (caller.APIID, error) {
	var id caller.APIID
	err := api.call("database", caller.EmptyParams, &id)
	return id, err
}

func (api *API) History() (caller.APIID, error) {
	var id caller.APIID
	err := api.call("history", caller.EmptyParams, &id)
	return id, err
}

func (api *API) NetworkBroadcast() (caller.APIID, error) {
	var id caller.APIID
	err := api.call("network_broadcast", caller.EmptyParams, &id)
	return id, err
}
