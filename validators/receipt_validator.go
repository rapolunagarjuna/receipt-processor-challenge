package validators

import (
	"github.com/go-playground/validator/v10"
	"time"
	"regexp"
)

func ValidateDate(fl validator.FieldLevel) bool {
	_, err := time.Parse("2006-01-02", fl.Field().String())
	return err == nil
}

func ValidateTime(fl validator.FieldLevel) bool {
	_, err := time.Parse("15:04", fl.Field().String())
	return err == nil
}

func ValidateName(fl validator.FieldLevel) bool {
	match, _ := regexp.MatchString("^[\\w\\s\\-&]+$", fl.Field().String())
	return match
}

func ValidateTotal(fl validator.FieldLevel) bool {
	match, _ := regexp.MatchString("^\\d+\\.\\d{2}$", fl.Field().String())
	return match
}