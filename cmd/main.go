package main

import (
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rapolunagarjuna/receipt-procesor-challenge/validators"
)

// Define the item struct with validation tags
type item struct {
	ShortDescription string `json:"shortDescription" validate:"required,name"`
	Price            string `json:"price" validate:"required,total"`
}

// Define the receipt struct with validation tags
type receipt struct {
	Retailer     string `json:"retailer" validate:"required,name"`
	PurchaseDate string `json:"purchaseDate" validate:"required,date"`
	PurchaseTime string `json:"purchaseTime" validate:"required,time"`
	Items        []item `json:"items" validate:"required,min=1,dive"`
	Total        string `json:"total" validate:"required,total"`
}

// Initialize the validator
var validate = validator.New()

// func validateDate(fl validator.FieldLevel) bool {
// 	_, err := time.Parse("2006-01-02", fl.Field().String())
// 	return err == nil
// }

// func validateTime(fl validator.FieldLevel) bool {
// 	_, err := time.Parse("15:04", fl.Field().String())
// 	return err == nil
// }

// func validateName(fl validator.FieldLevel) bool {
// 	match, _ := regexp.MatchString("^[\\w\\s\\-&]+$", fl.Field().String())
// 	return match
// }

// func validateTotal(fl validator.FieldLevel) bool {
// 	match, _ := regexp.MatchString("^\\d+\\.\\d{2}$", fl.Field().String())
// 	return match
// }

func init() {
	validate.RegisterValidation("date", ValidateDate)
	validate.RegisterValidation("time", ValidateTime)
	validate.RegisterValidation("total", ValidateTotal)
	validate.RegisterValidation("name", ValidateName)
}

func getTotalPoints(r *receipt) int64 {
	var points int64

	// One point for every alphanumeric character in the retailer name.
	for _, char := range r.Retailer {
		if ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || ('0' <= char && char <= '9') {
			points++
		}
	}

	// 50 points if the total is a round dollar amount with no cents.
	// 25 points if the total is a multiple of 0.25.
	if total, err := strconv.ParseFloat(r.Total, 64); err == nil {
		if total == math.Floor(total) {
			points += 50
		}
		if math.Mod(total, 0.25) == 0 {
			points += 25
		}
	}

	// 5 points for every two items on the receipt.
	points += int64(len(r.Items) / 2) * 5

	// If the trimmed length of the item description is a multiple of 3,
	// multiply the price by 0.2 and round up to the nearest integer.
	// The result is the number of points earned.
	for _, item := range r.Items {
		descriptionLenAfterTrim := len(strings.TrimSpace(item.ShortDescription))
		if descriptionLenAfterTrim%3 == 0 {
			if price, err := strconv.ParseFloat(item.Price, 64); err == nil {
				points += int64(math.Ceil(price * 0.2))
			}
		}
	}

	// 6 points if the day in the purchase date is odd.
	if date, err := time.Parse("2006-01-02", r.PurchaseDate); err == nil && date.Day()%2 != 0 {
		points += 6
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	if time, err := time.Parse("15:04", r.PurchaseTime); err == nil && time.Hour() >= 14 && time.Hour() < 16 {
		points += 10
	}

	return points
}

func main() {
	allReceipts := make(map[string]int64)

	r := gin.Default()
	log.Println("Starting server on port 8080")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/receipts/:id/points", func(c *gin.Context) {
		id := c.Param("id")

		points, ok := allReceipts[id]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"description": "No receipt found for that id"})
			return
		}

		c.JSON(200, gin.H{
			"points": points,
		})
	})

	r.POST("/receipts/process", func(c *gin.Context) {
		var newReceipt receipt

		if err := c.ShouldBindJSON(&newReceipt); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"description": "The receipt is invalid"})
			return
		}

		if err := validate.Struct(&newReceipt); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"description": "The receipt is invalid"})
			return
		}

		id := uuid.New().String()
		if _, ok := allReceipts[id]; ok {
			c.JSON(http.StatusInternalServerError, gin.H{"description": "Internal server error, Try again"})
			return
		}

		allReceipts[id] = getTotalPoints(&newReceipt)

		c.JSON(http.StatusOK, gin.H{
			"id": id,
		})
	})

	r.Run() // listens on 8080
}
