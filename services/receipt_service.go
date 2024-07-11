package services

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/rapolunagarjuna/receipt-processor-challenge/models"
	"github.com/rapolunagarjuna/receipt-processor-challenge/db"
)

/*
ReceiptService is an interface that contains the methods to interact with the receipt
AddNewReceipt is a method that adds a new receipt to the database
GetReceipt is a method that returns the points of the receipt
*/

type ReceiptService interface {
	AddNewReceipt(r *models.Receipt) (string, int64)
	GetReceipt(id string) (int64, bool)
}

/*
ReceiptServiceImpl is a struct that contains the DB
DB is an interface that contains the methods to interact with the database
AddNewReceipt is a method that adds a new receipt to the database
GetReceipt is a method that returns the points of the receipt

The reason for using an Implementation of the interface is to make the code more modular
and to make the code more testable
*/
type ReceiptServiceImpl struct {
	DB db.DB
}

/*
AddNewReceipt is a function that adds a new receipt to the database
it calculates the points of the receipt and adds the points to the database

assumption here is that the receipt is valid
and the conversions are successful
because the receipt is validated before calling this function
*/
func (receiptService *ReceiptServiceImpl) AddNewReceipt(r *models.Receipt) (string, int64) {
	var points int64

	points += PointsForRetailerName(r.Retailer)
	points += PointsForReceiptTotal(r.Total)
	points += PointsForItems(r.Items)
	points += PointsForItemDescription(r.Items)
	points += PointsForReceiptPurchaseDate(r.PurchaseDate)
	points += PointsForReceiptPurchaseTime(r.PurchaseTime)
	
	id := receiptService.DB.AddNewReceipt(points)
	return id, points
}

/*
GetReceipt is a function that returns the points of the receipt
if the receipt is not found, returns 404
*/
func (receiptService *ReceiptServiceImpl) GetReceipt(id string) (int64, bool) {
	if points, ok := receiptService.DB.GetReceipt(id); ok {
		return points, true
	}
	return int64(0), false
}

// One point for every alphanumeric character in the retailer name.
func PointsForRetailerName(retailerName string) int64 {
	var points int64

	for _, char := range retailerName {
		if ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || ('0' <= char && char <= '9') {
			points++
		}
	}
	return points
}

/*
 50 points if the total is a round dollar amount with no cents.
 25 points if the total is a multiple of 0.25.
 */
func PointsForReceiptTotal(receiptTotal string) int64 {
	var points int64
	if total, err := strconv.ParseFloat(receiptTotal, 64); err == nil {
		if total == math.Floor(total) {
			points += 50
		}
		if math.Mod(total, 0.25) == 0 {
			points += 25
		}
	}
	return points
}

// 5 points for every two items on the receipt.
func PointsForItems(items []models.Item) int64 {
	return int64(len(items)) / 2 * 5
}

/*
If the trimmed length of the item description is a multiple of 3,
multiply the price by 0.2 and round up to the nearest integer.
The result is the number of points earned.
*/
func PointsForItemDescription(items []models.Item) int64 {
	var points int64
	for _, item := range items {
		descriptionLenAfterTrim := len(strings.TrimSpace(item.ShortDescription))
		if descriptionLenAfterTrim%3 == 0 {
			if price, err := strconv.ParseFloat(item.Price, 64); err == nil {
				points += int64(math.Ceil(price * 0.2))
			}
		}
	}
	return points
}

// 6 points if the day in the purchase date is odd.
func PointsForReceiptPurchaseDate(purchaseDate string) int64 {
	var points int64
	if date, err := time.Parse("2006-01-02", purchaseDate); err == nil && date.Day()%2 != 0 {
		points += 6
	}
	return points
}

// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
func PointsForReceiptPurchaseTime(purchaseTime string) int64 {
	var points int64
	if time, err := time.Parse("15:04", purchaseTime); err == nil && time.Hour() >= 14 && time.Hour() < 16 {
		points += 10
	}
	return points
}



