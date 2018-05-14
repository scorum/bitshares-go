package types

import (
	"fmt"
	"github.com/pkg/errors"
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
