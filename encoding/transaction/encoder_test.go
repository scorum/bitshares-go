package transaction

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncoder_EncodeMoney_ZeroTrimming(t *testing.T) {
	buffer := new(bytes.Buffer)

	encoder := NewEncoder(buffer)
	err := encoder.EncodeMoney("00000000000000000000099.0000000000000000000000000 SCR")
	require.NoError(t, err)

	result := make([]byte, 16)
	buffer.Read(result)
	buffer.Reset()

	err = encoder.EncodeMoney("99.000 SCR")
	require.NoError(t, err)

	result2 := make([]byte, 16)
	buffer.Read(result2)

	require.Equal(t, result, result2)
}

func TestEncoder_EncodeMoney_ValueOverflow(t *testing.T) {
	buffer := new(bytes.Buffer)
	encoder := NewEncoder(buffer)

	err := encoder.EncodeMoney("11111111111111111111111111111111111111 SCR")
	require.Error(t, err)
}
