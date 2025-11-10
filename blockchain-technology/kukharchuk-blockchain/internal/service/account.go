package service

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"kukharchuk-blockchain/internal/model"
)

type Client interface {
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	NetworkID(ctx context.Context) (*big.Int, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
}

type AccountService struct {
	client Client
}

func NewAccountService(client Client) *AccountService {
	return &AccountService{
		client: client,
	}
}

func (s *AccountService) Balance(hexAddress string) (*big.Int, error) {
	var (
		ctx     = context.Background()
		address = common.HexToAddress(hexAddress)
	)

	balance, err := s.client.BalanceAt(ctx, address, nil)
	if err != nil {
		return nil, fmt.Errorf("get balance: %w", err)
	}

	return balance, nil
}

func (s *AccountService) SendTransaction(hexPrivateKey string, recipient string, amount *big.Int) (string, error) {
	var (
		ctx = context.Background()
	)

	privateKey, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return "", fmt.Errorf("convert: %w", err)
	}

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	nonce, err := s.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return "", fmt.Errorf("pending nonce: %w", err)
	}

	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("gas price: %w", err)
	}
	gasLimit := uint64(21000)

	toAddress := common.HexToAddress(recipient)
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &toAddress,
		Value:    amount,
	})

	chainID, err := s.client.NetworkID(ctx)
	if err != nil {
		return "", fmt.Errorf("get chained ID: %w", err)
	}

	signedTransaction, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", fmt.Errorf("sign transaction: %w", err)
	}

	err = s.client.SendTransaction(ctx, signedTransaction)
	if err != nil {
		return "", fmt.Errorf("send transaction: %w", err)
	}

	return signedTransaction.Hash().Hex(), nil
}

func (s *AccountService) AddAccount() (string, string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", fmt.Errorf("generate key: %w", err)
	}

	privateKeyBites := crypto.FromECDSA(privateKey)
	privateKeyHex := hexutil.Encode(privateKeyBites)[2:]

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", fmt.Errorf("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return address.Hex(), privateKeyHex, nil
}

func (s *AccountService) GetTransaction(hash string) (*model.TransactionData, error) {
	var (
		ctx    = context.Background()
		txHash = common.HexToHash(hash)
	)

	transaction, isPending, err := s.client.TransactionByHash(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("get transaction from ethereum: %v", err)
	}

	if isPending {
		return &model.TransactionData{
			Hash:   hash,
			Status: "pending",
		}, nil
	}

	receipt, err := s.client.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("get receipt by transaction hash: %v", err)
	}

	signer := types.NewEIP155Signer(transaction.ChainId())
	from, err := types.Sender(signer, transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender: %v", err)
	}

	valueETH := new(big.Float).Quo(
		new(big.Float).SetInt(transaction.Value()),
		big.NewFloat(1e18),
	).String()

	status := "failed"
	if receipt.Status == 1 {
		status = "success"
	}

	return &model.TransactionData{
		Hash:        hash,
		From:        from.Hex(),
		To:          transaction.To().Hex(),
		ValueETH:    valueETH,
		GasUsed:     receipt.GasUsed,
		BlockNumber: receipt.BlockNumber.Uint64(),
		Status:      status,
		Input:       fmt.Sprintf("0x%x", transaction.Data()),
	}, nil
}
