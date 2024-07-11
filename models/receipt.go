package models


// Receipt is a struct that contains the retailer, purchaseDate, purchaseTime, items and total of the receipt
type Receipt struct {
	Retailer     string `json:"retailer" validate:"required,alphanumeric"`
	PurchaseDate string `json:"purchaseDate" validate:"required,receiptDate"`
	PurchaseTime string `json:"purchaseTime" validate:"required,receiptTime"`
	Items        []Item `json:"items" validate:"required,min=1,dive"`
	Total        string `json:"total" validate:"required,decimal"`
}
