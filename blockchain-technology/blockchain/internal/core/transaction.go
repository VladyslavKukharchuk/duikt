package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Transaction struct {
	Sender    string
	Recipient string
	Amount    float64
}

func generateTransactionID(tx Transaction) string {
	data := fmt.Sprintf("%s%s%f", tx.Sender, tx.Recipient, tx.Amount)
	hash := sha256.Sum256([]byte(data))

	return hex.EncodeToString(hash[:])
}
