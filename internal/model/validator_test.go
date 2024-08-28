package model

import (
	"context"
	"testing"
	"time"
)

func TestOrderValid(t *testing.T) {
	tests := []struct {
		name     string
		order    Order
		expected map[string]string
	}{
		{
			name: "Test valid order",
			order: Order{
				UID:             "valid-uuid",
				TrackNumber:     "123456",
				Entry:           "entry",
				Items:           []ItemAttrs{{ChrtID: 1, Price: 100, Rid: "rid", Name: "name", Sale: 10, Size: "size", TotalPrice: 100, NmID: 1, Brand: "brand", Status: 1}},
				Locale:          "en",
				CustomerID:      "customer-id",
				DeliveryService: "delivery-service",
				ShardKey:        "shardkey",
				SmID:            1,
				DateCreated:     time.Now(),
				Delivery:        DeliveryAttrs{Name: "John Doe", Phone: "1234567890", Zip: "12345", City: "City", Address: "Address", Region: "Region", Email: "email@example.com"},
				Payment:         PaymentAttrs{Transaction: "valid-uuid", Currency: "USD", Provider: "provider", Amount: 100},
			},
			expected: map[string]string{},
		},
		{
			name: "Test invalid order with missing fields",
			order: Order{
				UID:             "",
				TrackNumber:     "",
				Entry:           "",
				Items:           []ItemAttrs{},
				Locale:          "",
				CustomerID:      "",
				DeliveryService: "",
				ShardKey:        "",
				SmID:            -1,
				DateCreated:     time.Time{},
				Delivery:        DeliveryAttrs{Name: "", Phone: "", Zip: "", City: "", Address: "", Region: "", Email: ""},
				Payment:         PaymentAttrs{Transaction: "", Currency: "", Provider: "", Amount: -1},
			},
			expected: map[string]string{
				"order_uid":        "Order UID is required and must be a valid UUID",
				"track_number":     "Track Number is required",
				"entry":            "Entry is required",
				"items":            "At least one item is required",
				"locale":           "Locale is required",
				"customer_id":      "Customer ID is required",
				"delivery_service": "Delivery Service is required",
				"shardkey":         "ShardKey is required",
				"sm_id":            "SmID must be a positive integer",
				"date_created":     "Date Created is required",
				"delivery":         "Delivery attributes are invalid",
				"payment":          "Payment attributes are invalid",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.order.Valid(context.Background())
			if len(got) != len(tt.expected) {
				t.Errorf("Expected %v problems, got %v", tt.expected, got)
			}
			for key, expectedMsg := range tt.expected {
				if gotMsg, ok := got[key]; !ok {
					t.Errorf("Missing expected problem for key %s", key)
				} else if gotMsg != expectedMsg {
					t.Errorf("For key %s, expected message %q, got %q", key, expectedMsg, gotMsg)
				}
			}
		})
	}
}

func TestDeliveryAttrsValid(t *testing.T) {
	tests := []struct {
		name     string
		delivery DeliveryAttrs
		expected map[string]string
	}{
		{
			name: "Test valid delivery",
			delivery: DeliveryAttrs{
				Name:    "John Doe",
				Phone:   "1234567890",
				Zip:     "12345",
				City:    "City",
				Address: "Address",
				Region:  "Region",
				Email:   "email@example.com",
			},
			expected: map[string]string{},
		},
		{
			name: "Test invalid delivery with missing fields",
			delivery: DeliveryAttrs{
				Name:    "",
				Phone:   "",
				Zip:     "",
				City:    "",
				Address: "",
				Region:  "",
				Email:   "",
			},
			expected: map[string]string{
				"name":    "Name is required",
				"phone":   "Phone is required",
				"zip":     "Zip is required",
				"city":    "City is required",
				"address": "Address is required",
				"region":  "Region is required",
				"email":   "Email is required",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.delivery.Valid(context.Background())
			if len(got) != len(tt.expected) {
				t.Errorf("Expected %v problems, got %v", tt.expected, got)
			}
			for key, expectedMsg := range tt.expected {
				if gotMsg, ok := got[key]; !ok {
					t.Errorf("Missing expected problem for key %s", key)
				} else if gotMsg != expectedMsg {
					t.Errorf("For key %s, expected message %q, got %q", key, expectedMsg, gotMsg)
				}
			}
		})
	}
}

func TestPaymentAttrs_Valid(t *testing.T) {
	tests := []struct {
		name     string
		payment  PaymentAttrs
		expected map[string]string
	}{
		{
			name: "Test valid payment",
			payment: PaymentAttrs{
				Transaction: "valid-uuid",
				Currency:    "USD",
				Provider:    "provider",
				Amount:      100,
			},
			expected: map[string]string{},
		},
		{
			name: "Test invalid payment with missing fields",
			payment: PaymentAttrs{
				Transaction: "",
				Currency:    "",
				Provider:    "",
				Amount:      -1,
			},
			expected: map[string]string{
				"transaction": "Transaction ID is required and must be a valid UUID",
				"currency":    "Currency is required",
				"provider":    "Provider is required",
				"amount":      "Amount must be a positive integer",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.payment.Valid(context.Background())
			if len(got) != len(tt.expected) {
				t.Errorf("Expected %v problems, got %v", tt.expected, got)
			}
			for key, expectedMsg := range tt.expected {
				if gotMsg, ok := got[key]; !ok {
					t.Errorf("Missing expected problem for key %s", key)
				} else if gotMsg != expectedMsg {
					t.Errorf("For key %s, expected message %q, got %q", key, expectedMsg, gotMsg)
				}
			}
		})
	}
}
