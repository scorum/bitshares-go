package types

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/scorum/openledger-go/encoding/transaction"
	"strconv"
	"strings"
)

// On the BitShares blockchains there are no addresses, but objects identified by a unique id,
// an type and a space in the form: space.type.id
type ObjectID struct {
	Space uint64
	Type  uint64
	ID    uint64
}

func (o ObjectID) String() string {
	return fmt.Sprintf("%d.%d.%d", o.Space, o.Type, o.ID)
}

func (o *ObjectID) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

func (o *ObjectID) UnmarshalJSON(s []byte) error {
	str, err := unquote(string(s))
	if err != nil {
		return errors.Errorf("unable to parse ObjectID from %s", s)
	}

	objectID, err := ParseObjectID(str)
	if err != nil {
		return err
	}

	*o = objectID
	return nil
}

func (o ObjectID) MarshalTransaction(encoder *transaction.Encoder) error {
	encoder.Encode(o.String())
	return nil
}

func MustParseObjectID(str string) ObjectID {
	out, err := ParseObjectID(str)
	if err != nil {
		panic(err)
	}
	return out
}

func ParseObjectID(str string) (ObjectID, error) {
	var err error

	objectID := ObjectID{}

	parts := strings.Split(str, ".")

	if len(parts) != 3 {
		return objectID, errors.Errorf("unable to parse ObjectID from %s", str)
	}

	objectID.Space, err = strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return objectID, errors.Errorf("unable to parse ObjectID [space] from %s", str)
	}

	objectID.Type, err = strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return objectID, errors.Errorf("unable to parse ObjectID [type] from %s", str)
	}

	objectID.ID, err = strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		return objectID, errors.Errorf("unable to parse ObjectID [id] from %ss", str)
	}

	return objectID, nil
}

func unquote(in string) (string, error) {
	if strings.HasPrefix(in, "\"") && strings.HasSuffix(in, "\"") {
		q, err := strconv.Unquote(in)
		if err != nil {
			return "", err
		}
		return q, nil
	}
	return in, nil
}
