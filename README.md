# scorum/scorum-go
[![Go Report Card](https://goreportcard.com/badge/github.com/scorum/openledger-go)](https://goreportcard.com/report/github.com/scorum/openledger-go)
[![GoDoc](https://godoc.org/github.com/scorum/openledger-go?status.svg)](https://godoc.org/github.com/scorum/openledger-go)

Golang RPC (via websockets) client library for [Bitshares](https://bitshares.org/) and [Openledger](https://openledger.io) in particular

## Usage

```go
import "github.com/scorum/openledger-go"
```

## Example
```go
client, _ := NewClient("wss://bitshares.openledger.info/ws")

// retrieve the current global_property_object
client.Database.GetDynamicGlobalProperties()

// set a block applied callback
client.Database.SetBlockAppliedCallback(func(blockID string, err error) {
    log.Println(blockID)
})

// cancel all callbacks
client.Database.CancelAllSubscriptions()

```