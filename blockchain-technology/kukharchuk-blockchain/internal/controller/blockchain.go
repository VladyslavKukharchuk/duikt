package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kukharchuk-blockchain/internal/model"
	"kukharchuk-blockchain/internal/service"
)

type BlockchainController struct {
	service *service.BlockchainService
}

func NewBlockchainController(service *service.BlockchainService) *BlockchainController {
	return &BlockchainController{service: service}
}

type BaseResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func newResponse(status string, data interface{}) BaseResponse {
	return BaseResponse{
		Status: status,
		Data:   data,
	}
}

func toBlocks(model []model.Block) []Block {
	blocks := make([]Block, 0, len(model))

	for _, b := range model {
		blocks = append(blocks, toBlock(b))
	}

	return blocks
}

type Block struct {
	Index        int           `json:"index"`
	Timestamp    int64         `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	Proof        int           `json:"proof"`
	PreviousHash string        `json:"previous_hash"`
	Hash         string        `json:"hash"`
}

func toBlock(model model.Block) Block {
	return Block{
		Index:        model.Index,
		Timestamp:    model.Timestamp,
		Transactions: toTransactions(model.Transactions),
		Proof:        model.Proof,
		PreviousHash: model.PreviousHash,
		Hash:         model.Hash,
	}
}

func toTransactions(model []model.Transaction) []Transaction {
	transactions := make([]Transaction, 0, len(model))

	for _, t := range model {
		transactions = append(transactions, toTransaction(t))
	}

	return transactions
}

type Transaction struct {
	Sender    string  `json:"sender"`
	Recipient string  `json:"recipient"`
	Amount    float64 `json:"amount"`
}

func toTransaction(model model.Transaction) Transaction {
	return Transaction{
		Sender:    model.Sender,
		Recipient: model.Recipient,
		Amount:    model.Amount,
	}
}

func (c *BlockchainController) GetBlockchain(ctx *gin.Context) {
	blockchain := c.service.GetChain()

	res := newResponse("success", toBlocks(blockchain))
	ctx.JSON(http.StatusOK, res)
}

type TransactionRequest struct {
	Sender    string  `json:"sender"`
	Recipient string  `json:"recipient"`
	Amount    float64 `json:"amount"`
}

func (t TransactionRequest) toTransactionModel() model.Transaction {
	return model.Transaction{
		Sender:    t.Sender,
		Recipient: t.Recipient,
		Amount:    t.Amount,
	}
}

type TransactionResponse struct {
	TransactionID string `json:"transaction_id"`
}

func (c *BlockchainController) AddTransaction(ctx *gin.Context) {
	var transaction TransactionRequest
	if err := ctx.ShouldBindJSON(&transaction); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponse("error", "Invalid transaction data"))
		return
	}

	transactionID := c.service.AddTransaction(transaction.toTransactionModel())

	res := newResponse("success", TransactionResponse{TransactionID: transactionID})
	ctx.JSON(http.StatusOK, res)
}

type BlockResponse struct {
	Index        int    `json:"index"`
	ProofOfWork  int    `json:"proof_of_work"`
	Timestamp    int64  `json:"timestamp"`
	PreviousHash string `json:"previous_hash"`
}

func (c *BlockchainController) MineBlock(ctx *gin.Context) {
	index, proof, candidateTimestamp, previousHash := c.service.MineBlock()

	res := newResponse("success", BlockResponse{
		Index:        index,
		ProofOfWork:  proof,
		Timestamp:    candidateTimestamp,
		PreviousHash: previousHash,
	})
	ctx.JSON(http.StatusOK, res)
}
