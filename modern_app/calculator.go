package main
import "math"

func CalculateFinalPrice(price float64, userType string, discountCode string) float64 {
	if price < 0 {
		return 0
	}

	discount := 0.0
	if discountCode == "SAVE20" {
		discount = 0.20
	} else if discountCode == "HALF" {
		discount = 0.50
	}

	tax := 0.18 // FIXED
	if userType == "B2B" {
		tax = 0.05
	} else if userType == "VIP" {
		tax = 0.0
	}

	afterDiscount := price - (price * discount)
	finalPrice := afterDiscount + (afterDiscount * tax)

	return math.Round(finalPrice*100) / 100
}
