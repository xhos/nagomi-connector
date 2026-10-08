package wise

import (
	"testing"

	wiseapi "nagomi-connector/internal/wise"
)

func TestToDomainDescription(t *testing.T) {
	tests := []struct {
		name    string
		details wiseapi.Details
		want    string
	}{
		{
			name: "card payment uses merchant and city",
			details: wiseapi.Details{
				Type:        "CARD",
				Description: "Card transaction of 237 JPY issued by Seven-Eleven TOKYO",
				Merchant:    &wiseapi.Merchant{Name: "Seven-Eleven", City: "TOKYO"},
			},
			want: "Seven-Eleven TOKYO",
		},
		{
			name: "merchant without city",
			details: wiseapi.Details{
				Type:        "CARD",
				Description: "Card transaction of 9.99 USD issued by Netflix",
				Merchant:    &wiseapi.Merchant{Name: "Netflix"},
			},
			want: "Netflix",
		},
		{
			name: "no merchant keeps wise's description",
			details: wiseapi.Details{
				Type:        "CONVERSION",
				Description: "Converted 100.00 CAD to 10,800 JPY",
			},
			want: "Converted 100.00 CAD to 10,800 JPY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toDomain(&wiseapi.Transaction{Type: "DEBIT", Details: tt.details}, 1)
			if got.Description != tt.want {
				t.Errorf("Description = %q; want %q", got.Description, tt.want)
			}
		})
	}
}
