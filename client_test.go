package openledger

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClient(t *testing.T) {
	t.Run("valid ws url", func(t *testing.T) {
		const url = "wss://bitshares.openledger.info/ws"
		_, err := NewClient(url)
		require.NoError(t, err)
	})

	t.Run("invalid ws url", func(t *testing.T) {
		_, err := NewClient("wss://invalid")
		require.Error(t, err)
	})
}
