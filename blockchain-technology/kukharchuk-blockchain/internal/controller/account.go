package controller

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

type AccountService interface {
	Balance(hexAddress string) (*big.Int, error)
	SendTransaction(hexPrivateKey string, recipient string, amount *big.Int) (string, error)
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
