package controller

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

type AccountService interface {
	Balance(xehAddress string) (*big.Int, error)
}

type AccountController struct {
	service AccountService
}

func NewAccountController(service AccountService) *AccountController {
	return &AccountController{service: service}
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

	floatBalance := new(big.Float).Quo(
		new(big.Float).SetInt(balance),
		big.NewFloat(1e18),
	)

	res := newResponse("success", floatBalance)
	ctx.JSON(http.StatusOK, res)
}
