package main

import (
	"testing"
)

func TestCalculateFinalPrice(t *testing.T) {
	tests := []struct {
		name         string
		price        float64
		userType     string
		discountCode string
		expected     float64
	}{
		{"Negative price", -10.0, "NORMAL", "NONE", 0.0},
		{"Normal user, no discount", 100.0, "NORMAL", "NONE", 118.0}, // 100 + 18%
		{"B2B user, SAVE20", 100.0, "B2B", "SAVE20", 84.0},         // 100 - 20% = 80 + 5% = 84
		{"VIP user, HALF", 100.0, "VIP", "HALF", 50.0},             // 100 - 50% = 50 + 0% = 50
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateFinalPrice(tt.price, tt.userType, tt.discountCode)
			if got != tt.expected {
				t.Errorf("CalculateFinalPrice() = %v, want %v", got, tt.expected)
			}
		})
	}
}
