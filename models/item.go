package models

// Item is a struct that contains the shortDescription and the price of the item

type Item struct {
	ShortDescription string `json:"shortDescription" validate:"required,alphanumeric"`
	Price            string `json:"price" validate:"required,decimal"`
}