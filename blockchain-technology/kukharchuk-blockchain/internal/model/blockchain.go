package model

import (
	"strings"
	"time"
)

type Blockchain struct {
	Chain   []Block
	Mempool []Transaction
}

func NewBlockchain() *Blockchain {
	bc := &Blockchain{
		Chain:   []Block{},
		Mempool: []Transaction{},
	}

	bc.createGenesisBlock()

	return bc
}

func (bc *Blockchain) createGenesisBlock() {
	genesisBlock := Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Transactions: []Transaction{},
		Proof:        11122000,
		PreviousHash: "kukharchuk",
	}

	genesisBlock.Hash = calculateHash(genesisBlock)

	bc.Chain = append(bc.Chain, genesisBlock)
}

func (bc *Blockchain) AddBlock(proof int, timestamp int64, previousHash string) {
	newBlock := Block{
		Index:        len(bc.Chain),
		Timestamp:    timestamp,
		Transactions: bc.Mempool,
		Proof:        proof,
		PreviousHash: previousHash,
	}

	newBlock.Hash = calculateHash(newBlock)

	bc.Chain = append(bc.Chain, newBlock)

	bc.Mempool = []Transaction{}
}

func (bc *Blockchain) AddTransaction(transaction Transaction) string {
	bc.Mempool = append(bc.Mempool, transaction)

	return generateTransactionID(transaction)
}

func (bc *Blockchain) ProofOfWork() (int, int64) {
	lastBlock := bc.Chain[len(bc.Chain)-1]
	candidateTimestamp := time.Now().Unix()
	proof := 0

	for !validProof(lastBlock, proof, bc.Mempool, candidateTimestamp) {
		proof++
	}

	return proof, candidateTimestamp
}

func validProof(lastBlock Block, proof int, transactions []Transaction, candidateTimestamp int64) bool {
	tempBlock := Block{
		Index:        lastBlock.Index + 1,
		Timestamp:    candidateTimestamp,
		Transactions: transactions,
		Proof:        proof,
		PreviousHash: lastBlock.Hash,
	}

	guessHash := calculateHash(tempBlock)

	monthSuffix := "12"

	return guessHash[:4] == "0000" && strings.HasSuffix(guessHash, monthSuffix)
}
