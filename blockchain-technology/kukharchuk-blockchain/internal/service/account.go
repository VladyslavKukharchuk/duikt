package service

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Client interface {
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
}

type AccountService struct {
	client Client
}

func NewAccountService(client Client) *AccountService {
	return &AccountService{
		client: client,
	}
}

func (s *AccountService) Balance(xehAddress string) (*big.Int, error) {
	var (
		ctx     = context.Background()
		address = common.HexToAddress(xehAddress)
	)

	balance, err := s.client.BalanceAt(ctx, address, nil)
	if err != nil {
		return nil, fmt.Errorf("get balance: %w", err)
	}

	return balance, nil
}
