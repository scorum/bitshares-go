package types

import (
	"encoding/json"
	"strconv"
)

// Suint64 uint64 with redeclared JSON unmarshal;
// Can be parsed from uint64 either string
type Suint64 uint64

// Suint32 uint32 with redeclared JSON unmarshal;
// Can be parsed from uint32 either string
type Suint32 uint32

func (su *Suint64) UnmarshalJSON(b []byte) (err error) {
	var u uint64
	if err = json.Unmarshal(b, &u); err == nil {
		temp := Suint64(u)
		su = &temp
		return nil
	}

	// failed on uint64, try string
	var s string
	if err = json.Unmarshal(b, &s); err == nil {
		u, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		temp := Suint64(u)
		su = &temp
		return nil
	}

	return err
}

func (su *Suint32) UnmarshalJSON(b []byte) (err error) {
	var u uint32
	if err = json.Unmarshal(b, &u); err == nil {
		temp := Suint32(u)
		su = &temp
		return nil
	}

	// failed on uint32, try string
	var s string
	if err = json.Unmarshal(b, &s); err == nil {
		u, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		temp := Suint32(u)
		su = &temp
		return nil
	}

	return err
}
