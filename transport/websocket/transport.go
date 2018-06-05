package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/scorum/openledger-go/caller"
	"log"
	"math"
	"strconv"
	"sync"

	"github.com/pkg/errors"
	"github.com/scorum/scorum-go/transport"
	"golang.org/x/net/websocket"
)

type Transport struct {
	conn *websocket.Conn

	reqMutex  sync.Mutex
	requestID uint64
	pending   map[uint64]*callRequest

	callbackMutex sync.Mutex
	callbackID    uint64
	callbacks     map[uint64]func(args json.RawMessage)

	closing  bool // user has called Close
	shutdown bool // server has told us to stop

	mutex sync.Mutex
}

// Represent an async call
type callRequest struct {
	Error error            // after completion, the error status.
	Done  chan bool        // strobes when call is complete.
	Reply *json.RawMessage // reply message
}

func NewTransport(url string) (*Transport, error) {
	ws, err := websocket.Dial(url, "", "http://localhost")
	if err != nil {
		return nil, err
	}

	client := &Transport{
		conn:      ws,
		pending:   make(map[uint64]*callRequest),
		callbacks: make(map[uint64]func(args json.RawMessage)),
	}

	go client.input()
	return client, nil
}

func (caller *Transport) Call(api caller.APIID, method string, args []interface{}, reply interface{}) error {
	caller.reqMutex.Lock()
	defer caller.reqMutex.Unlock()

	caller.mutex.Lock()
	if caller.closing || caller.shutdown {
		caller.mutex.Unlock()
		return transport.ErrShutdown
	}

	// increase request id
	if caller.requestID == math.MaxUint64 {
		caller.requestID = 0
	}
	caller.requestID++
	seq := caller.requestID

	c := &callRequest{
		Done: make(chan bool, 1),
	}
	caller.pending[seq] = c
	caller.mutex.Unlock()

	request := transport.RPCRequest{
		Method: "call",
		ID:     caller.requestID,
		Params: []interface{}{api, method, args},
	}

	debug, _ := json.Marshal(request)
	log.Printf("[DEBUG] ws.Call method:%s request:%s\n", method, string(debug))

	// send Json Rcp request
	if err := websocket.JSON.Send(caller.conn, request); err != nil {
		caller.mutex.Lock()
		delete(caller.pending, seq)
		caller.mutex.Unlock()
		return err
	}

	// wait for the call to complete
	<-c.Done
	if c.Error != nil {
		return c.Error
	}

	if c.Reply != nil {
		if err := json.Unmarshal(*c.Reply, reply); err != nil {
			return err
		}
	}
	return nil
}

func (caller *Transport) input() {
	for {
		var message string
		if err := websocket.Message.Receive(caller.conn, &message); err != nil {
			caller.stop(err)
			return
		}

		var response transport.RPCResponse
		if err := json.Unmarshal([]byte(message), &response); err != nil {
			caller.stop(err)
			return
		} else {
			if call, ok := caller.pending[response.ID]; ok {
				caller.onCallResponse(response, call)
			} else {
				//the message is not a pending call, but probably a callback notice
				var incoming transport.RPCIncoming
				if err := json.Unmarshal([]byte(message), &incoming); err != nil {
					caller.stop(err)
					return
				}
				if incoming.Method == "notice" {
					if err := caller.onNotice(incoming); err != nil {
						caller.stop(err)
						return
					}
				} else {
					log.Printf("protocol error: unknown message received: %+v\n", incoming)
				}
			}
		}
	}
}

// Return pending clients and shutdown the client
func (caller *Transport) stop(err error) {
	caller.reqMutex.Lock()
	caller.shutdown = true
	for _, call := range caller.pending {
		call.Error = err
		call.Done <- true
	}
	caller.reqMutex.Unlock()
}

// Call response handler
func (caller *Transport) onCallResponse(response transport.RPCResponse, call *callRequest) {
	caller.mutex.Lock()
	delete(caller.pending, response.ID)
	if response.Error != nil {
		call.Error = response.Error
	}
	call.Reply = response.Result
	call.Done <- true
	caller.mutex.Unlock()

	debug, _ := json.Marshal(response)
	println(string(debug))
}

// Incoming notice handler
func (caller *Transport) onNotice(incoming transport.RPCIncoming) error {
	length := len(incoming.Params)

	if length == 0 {
		return nil
	}

	if length == 1 {
		return fmt.Errorf("invalid notice(%+v) message with odd number of params", incoming)
	}

	for i := 0; i < length; i += 2 {
		callbackID, err := strconv.ParseUint(string(incoming.Params[i]), 10, 64)
		if err != nil {
			return errors.Wrapf(err, "failed to parse %s as callbackID in notice %+v", incoming.Params[i], incoming)
		}

		notice := caller.callbacks[callbackID]
		if notice == nil {
			return fmt.Errorf("callback %d is not registered", callbackID)
		}

		// invoke callback
		notice(incoming.Params[i+1])
	}

	return nil
}

func (caller *Transport) SetCallback(api caller.APIID, method string, notice func(args json.RawMessage)) error {
	// increase callback id
	caller.callbackMutex.Lock()
	if caller.callbackID == math.MaxUint64 {
		caller.callbackID = 0
	}
	caller.callbackID++
	caller.callbacks[caller.callbackID] = notice
	caller.callbackMutex.Unlock()

	return caller.Call(api, method, []interface{}{caller.callbackID}, nil)
}

// Close calls the underlying web socket Close method. If the connection is already
// shutting down, ErrShutdown is returned.
func (caller *Transport) Close() error {
	caller.mutex.Lock()
	if caller.closing {
		caller.mutex.Unlock()
		return transport.ErrShutdown
	}
	caller.closing = true
	caller.mutex.Unlock()
	return caller.conn.Close()
}
