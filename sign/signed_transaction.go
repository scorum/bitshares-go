package sign

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/scorum/openledger-go/encoding/transaction"
	"github.com/scorum/openledger-go/types"
	"log"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

type SignedTransaction struct {
	*types.Transaction
}

func NewSignedTransaction(tx *types.Transaction) *SignedTransaction {
	return &SignedTransaction{tx}
}

func (tx *SignedTransaction) Serialize() ([]byte, error) {
	var b bytes.Buffer
	encoder := transaction.NewEncoder(&b)

	if err := encoder.Encode(tx.Transaction); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (tx *SignedTransaction) Digest(chain string) ([]byte, error) {
	var msgBuffer bytes.Buffer

	// Write the chain ID.
	rawChainID, err := hex.DecodeString(chain)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode chain ID: %v", chain)
	}

	if _, err := msgBuffer.Write(rawChainID); err != nil {
		return nil, errors.Wrap(err, "failed to write chain ID")
	}

	// Write the serialized transaction.
	rawTx, err := tx.Serialize()
	if err != nil {
		return nil, err
	}

	if _, err := msgBuffer.Write(rawTx); err != nil {
		return nil, errors.Wrap(err, "failed to write serialized transaction")
	}

	msgBytes := msgBuffer.Bytes()
	message := hex.EncodeToString(msgBytes)
	log.Printf("[DEBUG] Digest final message:%s\n", message)

	// Compute the digest.
	digest := sha256.Sum256(msgBytes)
	return digest[:], nil
}

func (tx *SignedTransaction) Sign(wifs []string, chain string) error {
	digest, err := tx.Digest(chain)
	if err != nil {
		return err
	}

	privKeys := make([]*btcec.PrivateKey, len(wifs))
	for index, wif := range wifs {
		w, err := btcutil.DecodeWIF(wif)
		if err != nil {
			return err
		}
		privKeys[index] = w.PrivKey
	}

	// Set the signature array in the transaction.
	sigsHex := make([]string, len(privKeys))
	for index, privKey := range privKeys {
		sig := SignBufferSha256(digest, privKey.ToECDSA())
		sigsHex[index] = hex.EncodeToString(sig)
	}
	tx.Transaction.Signatures = sigsHex
	return nil
}
