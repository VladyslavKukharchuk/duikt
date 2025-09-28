package main

import (
	"fmt"
	"time"

	"kukharchuk-blockchain/internal/core"
)

func main() {
	bc := core.NewBlockchain()
	bc.AddTransaction("Alice", "Bob", 50)
	bc.AddTransaction("Bob", "Charlie", 25)

	start := time.Now()

	proof, candidateTimestamp := bc.ProofOfWork()

	duration := time.Since(start)
	fmt.Printf("Proof of Work знайдено: %d (час виконання: %s)\n", proof, duration)

	previousHash := bc.Chain[len(bc.Chain)-1].Hash
	bc.AddBlock(proof, candidateTimestamp, previousHash)

	fmt.Println("Blockchain:", bc.Chain)
}
