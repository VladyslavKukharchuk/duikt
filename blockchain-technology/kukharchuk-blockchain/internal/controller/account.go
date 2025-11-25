package controller

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"

	"kukharchuk-blockchain/internal/model"
)

type AccountService interface {
	Balance(hexAddress string) (*big.Int, error)
	SendTransaction(hexPrivateKey string, recipient string, amount *big.Int) (string, error)
	AddAccount() (string, string, error)
	GetTransaction(hash string) (*model.TransactionData, error)
}

type AccountController struct {
	service AccountService
}

func NewAccountController(service AccountService) *AccountController {
	return &AccountController{service: service}
}

type BalanceResponse struct {
	Address string `json:"address"`
	Wei     string `json:"wei"`
	Eth     string `json:"eth"`
}

func (c *AccountController) Balance(ctx *gin.Context) {
	address := ctx.Param("address")
	if !common.IsHexAddress(address) {
		ctx.JSON(http.StatusBadRequest, newResponse("error", fmt.Sprintf("invalid Ethereum address: %s", address)))
		return
	}

	balance, err := c.service.Balance(address)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponse("error", fmt.Sprintf("failed to get balance for address %s, error: %v", address, err)))
		return
	}

	eth := new(big.Float).Quo(
		new(big.Float).SetInt(balance),
		big.NewFloat(1e18),
	)

	response := BalanceResponse{
		Address: address,
		Wei:     balance.String(),
		Eth:     eth.Text('f', 18),
	}

	res := newResponse("success", response)
	ctx.JSON(http.StatusOK, res)
}

type SendTransactionRequest struct {
	SenderPrivateKey string `json:"sender_private_key"`
	Recipient        string `json:"recipient"`
	Amount           string `json:"amount"`
}

type SendTransactionResponse struct {
	TransactionHex string `json:"transaction_hex"`
}

func (c *AccountController) SendTransaction(ctx *gin.Context) {
	var transaction SendTransactionRequest
	if err := ctx.ShouldBindJSON(&transaction); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponse("error", "Invalid transaction data"))
		return
	}

	amount := new(big.Int)
	if _, ok := amount.SetString(transaction.Amount, 10); !ok {
		ctx.JSON(http.StatusBadRequest, newResponse("error", "invalid amount value"))
		return
	}

	transactionHex, err := c.service.SendTransaction(transaction.SenderPrivateKey, transaction.Recipient, amount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponse("error", fmt.Sprintf("failed to send transaction: %v", err)))
		return
	}

	res := newResponse("success", SendTransactionResponse{transactionHex})
	ctx.JSON(http.StatusOK, res)
}

type AddAccountResponse struct {
	Address       string `json:"address"`
	PrivateKeyHex string `json:"private_key_hex"`
}

func (c *AccountController) AddAccount(ctx *gin.Context) {
	address, privateKeyHex, err := c.service.AddAccount()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponse("error", fmt.Sprintf("failed to add account: %v", err)))
		return
	}

	res := newResponse("success", AddAccountResponse{
		Address:       address,
		PrivateKeyHex: privateKeyHex,
	})
	ctx.JSON(http.StatusOK, res)
}

type GetTransactionResponse struct {
	Hash        string `json:"hash"`
	From        string `json:"from"`
	To          string `json:"to"`
	ValueETH    string `json:"value_eth"`
	GasUsed     uint64 `json:"gas_used"`
	BlockNumber uint64 `json:"block_number"`
	Status      string `json:"status"`
	Input       string `json:"input"`
}

func toGetTransactionResponse(m *model.TransactionData) GetTransactionResponse {
	return GetTransactionResponse{
		Hash:        m.Hash,
		From:        m.From,
		To:          m.To,
		ValueETH:    m.ValueETH,
		GasUsed:     m.GasUsed,
		BlockNumber: m.BlockNumber,
		Status:      m.Status,
		Input:       m.Input,
	}
}

func (c *AccountController) GetTransaction(ctx *gin.Context) {
	hash := ctx.Param("hash")

	transaction, err := c.service.GetTransaction(hash)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponse("error", fmt.Errorf("failed to get transaction by hash: %s, error: %v", hash, err)))
		return
	}

	ctx.JSON(http.StatusOK, newResponse("success", toGetTransactionResponse(transaction)))
}
