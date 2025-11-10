package main

import (
	"github.com/gin-gonic/gin"

	"kukharchuk-blockchain/internal/controller"
	"kukharchuk-blockchain/internal/service"
)

func main() {
	blockchainService := service.NewBlockchainService()
	blockchainController := controller.NewBlockchainController(blockchainService)

	r := gin.Default()

	r.GET("/blockchain", blockchainController.GetBlockchain)
	r.POST("/transaction", blockchainController.AddTransaction)
	r.GET("/mine", blockchainController.MineBlock)

	r.Run(":8088")
}
