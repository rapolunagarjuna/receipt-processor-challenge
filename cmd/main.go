package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rapolunagarjuna/receipt-processor-challenge/services"
	"github.com/rapolunagarjuna/receipt-processor-challenge/controllers"
	"github.com/rapolunagarjuna/receipt-processor-challenge/db"
)

// server, database, receiptService, receiptController are the global variables
var (
	server = gin.Default()
	database = db.InMemoryDB{AllReceipts: make(map[string]int64)}
	receiptService = services.ReceiptServiceImpl{DB: &database}
	receiptController = controllers.ReceiptController{ReceiptService: &receiptService}
)

func main() {

	/*
	creating a group for all the receipt related routes /receipts endpoints
	consists of the following endpoints:
	1. GET /receipts/:id/points         -> returns the points for a given receipt id, 
											if the receipt is not found, returns 404
	2. POST /receipts/process			-> processes the receipt and returns the id of the receipt, 
											if the receipt is invalid, returns 400
	*/
	receiptApiRoutes := server.Group("/receipts") 
	{
		receiptApiRoutes.GET("/:id/points", receiptController.GetReceiptPoints)
		receiptApiRoutes.POST("/process", receiptController.ProcessReceipt)
	}
	
	server.Run(":8080")
}
