package service

import "kukharchuk-blockchain/internal/model"

type BlockchainService struct {
	bc *model.Blockchain
}

func NewBlockchainService() *BlockchainService {
	return &BlockchainService{
		bc: model.NewBlockchain(),
	}
}

func (s *BlockchainService) GetChain() []model.Block {
	return s.bc.Chain
}

func (s *BlockchainService) AddTransaction(transaction model.Transaction) string {
	return s.bc.AddTransaction(transaction)
}

func (s *BlockchainService) MineBlock() (int, int, int64, string) {
	proof, candidateTimestamp := s.bc.ProofOfWork()

	prevHash := s.bc.Chain[len(s.bc.Chain)-1].Hash
	s.bc.AddBlock(proof, candidateTimestamp, prevHash)

	index := len(s.bc.Chain) - 1

	return index, proof, candidateTimestamp, prevHash
}
