package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Block struct {
	Index        int
	Timestamp    int64
	Transactions []Transaction
	Proof        int
	PreviousHash string
	Hash         string
}

func calculateHash(block Block) string {
	hashInput := fmt.Sprintf("%d%d%d%s%d",
		block.Index,
		block.Timestamp,
		block.Proof,
		block.PreviousHash,
		len(block.Transactions))

	for _, tx := range block.Transactions {
		hashInput += tx.Sender + tx.Recipient + fmt.Sprintf("%f", tx.Amount)
	}

	hash := sha256.Sum256([]byte(hashInput))

	return hex.EncodeToString(hash[:])
}
