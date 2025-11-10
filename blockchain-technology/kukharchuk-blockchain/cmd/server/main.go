package main

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"

	"kukharchuk-blockchain/internal/controller"
	"kukharchuk-blockchain/internal/service"
)

func main() {
	blockchainService := service.NewBlockchainService()
	blockchainController := controller.NewBlockchainController(blockchainService)

	client, err := ethclient.Dial("https://ethereum-sepolia-rpc.publicnode.com")
	if err != nil {
		log.Fatal(err)
	}
	accountService := service.NewAccountService(client)
	accountController := controller.NewAccountController(accountService)

	r := gin.Default()

	r.GET("/blockchain", blockchainController.GetBlockchain)
	r.POST("/transaction", blockchainController.AddTransaction)
	r.GET("/mine", blockchainController.MineBlock)
	r.GET("/account/:address/balance", accountController.Balance)
	r.POST("/account/transaction", accountController.SendTransaction)

	r.Run(":8088")
}
