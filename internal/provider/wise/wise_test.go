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

func TestToDomainTransferRef(t *testing.T) {
	w := &wiseapi.Transaction{
		Type:            "DEBIT",
		ReferenceNumber: "BALANCE-5104068188",
		Details:         wiseapi.Details{Type: "CONVERSION", Description: "Converted 1,000.00 CAD to 114,541 JPY"},
	}
	got := toDomain(w, 1)
	if got.TransferRef != "wise:BALANCE-5104068188" {
		t.Errorf("TransferRef = %q", got.TransferRef)
	}
	if got.ExternalID != "BALANCE-5104068188" {
		t.Errorf("ExternalID = %q; the dedup key must not change", got.ExternalID)
	}
}
