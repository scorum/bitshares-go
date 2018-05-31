package sign

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/scorum/scorum-go/types"
	"github.com/stretchr/testify/require"
)

func TestTransaction_Digest(t *testing.T) {
	var tx *types.Transaction
	// Prepare the transaction.
	expiration := time.Date(2016, 8, 8, 12, 24, 17, 0, time.UTC)
	tx = &types.Transaction{
		RefBlockNum:    36029,
		RefBlockPrefix: 1164960351,
		Expiration:     &types.Time{&expiration},
	}
	tx.PushOperation(&types.VoteOperation{
		Voter:    "xeroc",
		Author:   "xeroc",
		Permlink: "piston",
		Weight:   10000,
	})

	const initChain = "0000000000000000000000000000000000000000000000000000000000000000"

	expected := "582176b1daf89984bc8b4fdcb24ff1433d1eb114a8c4bf20fb22ad580d035889"
	stx := NewSignedTransaction(tx)
	digest, err := stx.Digest(initChain)
	require.NoError(t, err)
	require.Equal(t, expected, hex.EncodeToString(digest))
}
