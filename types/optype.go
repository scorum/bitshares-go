package types

type OpType uint16

const (
	TransferOpType OpType = iota
	LimitOrderCreateOpType
	LimitOrderCancelOpType
	CallOrderUpdate
)
