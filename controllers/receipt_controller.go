package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rapolunagarjuna/receipt-processor-challenge/services"
	"github.com/rapolunagarjuna/receipt-processor-challenge/models"
	"github.com/go-playground/validator/v10"
	"github.com/rapolunagarjuna/receipt-processor-challenge/validators"
	"net/http"
)
/*
ReceiptController is a struct that contains the ReceiptService
perfoming dependency injection on the ReceiptService
*/
type ReceiptController struct {
	ReceiptService services.ReceiptService
}

/*
ProcessReceipt is a function that processes the receipt and returns the id of the receipt
if the receipt is invalid, returns 400
making use of the validator to validate the receipt
validating the 
purchaseDate            		-> must be present and should be a valid date format (YYYY-MM-DD)
purchaseTime 		  			-> must be present and should be a valid time format (HH:MM:SS)	
retailer 			  			-> must be present and should be a valid name of the form ^[\\w\\s\\-&]+$
total 			  				-> must be present and should be a valid total of the form ^\\d+\\.\\d{2}$
items 			  				-> must have atleast one item and should be a valid array of items
		shortDescription		-> must be present and should be a valid name of the form ^[\\w\\s\\-&]+$
		price					-> must be present and should be a valid price of the form ^\\d+\\.\\d{2}$
*/
func (controller *ReceiptController) ProcessReceipt(c *gin.Context) {
	var validate = validator.New()
	validate.RegisterValidation("receiptDate", validators.ValidateReceiptDate)
	validate.RegisterValidation("receiptTime", validators.ValidateReceiptTime)
	validate.RegisterValidation("decimal", validators.ValidateDecimal)
	validate.RegisterValidation("alphanumeric", validators.ValidateAlphanumeric)
	var newReceipt models.Receipt

	if err := c.ShouldBindJSON(&newReceipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": "The receipt is invalid"})
		return
	}

	if err := validate.Struct(&newReceipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": "The receipt is invalid"})
		return
	}
	
	id,_ := controller.ReceiptService.AddNewReceipt(&newReceipt)

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

/*
GetReceiptPoints is a function that returns the points of the receipt
if the receipt is not found, returns 404
*/
func (controller *ReceiptController) GetReceiptPoints(c *gin.Context) {
	id := c.Param("id")

	points, ok := controller.ReceiptService.GetReceipt(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"description": "No receipt found for that id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"points": points,
	})
}