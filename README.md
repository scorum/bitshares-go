# scorum/openledger-go
[![Go Report Card](https://goreportcard.com/badge/github.com/scorum/openledger-go)](https://goreportcard.com/report/github.com/scorum/openledger-go)
[![GoDoc](https://godoc.org/github.com/scorum/openledger-go?status.svg)](https://godoc.org/github.com/scorum/openledger-go)

Golang RPC (via websockets) client library for [Bitshares](https://bitshares.org/) and [OpenLedger](https://openledger.io) in particular

## Usage

```go
import "github.com/scorum/openledger-go"
```

## Example
```go
client, err := NewClient("wss://bitshares.openledger.info/ws")

// retrieve the current global_property_object
client.Database.GetDynamicGlobalProperties()

// lookup symbols ids
symbols, err := client.Database.LookupAssetSymbols("OPEN.SCR", "USD")
require.NoError(t, err)

openSCR := symbols[0].ID
USD := symbols[1].ID

// retrieve 5 last filled orders
orders, err := client.History.GetFillOrderHistory(openSCR, USD, 5)

// set a block applied callback
client.Database.SetBlockAppliedCallback(func(blockID string, err error) {
    log.Println(blockID)
})

// cancel all callbacks
client.Database.CancelAllSubscriptions()

```