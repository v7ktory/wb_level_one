package model

import (
	"context"
)

type Validator interface {
	Valid(ctx context.Context) map[string]string
}

func (o Order) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if o.UID == "" {
		problems["order_uid"] = "Order UID is required and must be a valid UUID"
	}
	if o.TrackNumber == "" {
		problems["track_number"] = "Track Number is required"
	}
	if o.Entry == "" {
		problems["entry"] = "Entry is required"
	}
	if len(o.Items) == 0 {
		problems["items"] = "At least one item is required"
	}
	if o.Locale == "" {
		problems["locale"] = "Locale is required"
	}
	if o.CustomerID == "" {
		problems["customer_id"] = "Customer ID is required"
	}
	if o.DeliveryService == "" {
		problems["delivery_service"] = "Delivery Service is required"
	}
	if o.ShardKey == "" {
		problems["shardkey"] = "ShardKey is required"
	}
	if o.SmID <= 0 {
		problems["sm_id"] = "SmID must be a positive integer"
	}
	if o.DateCreated.IsZero() {
		problems["date_created"] = "Date Created is required"
	}

	if p := o.Delivery.Valid(ctx); len(p) > 0 {
		problems["delivery"] = "Delivery attributes are invalid"
	}
	if p := o.Payment.Valid(ctx); len(p) > 0 {
		problems["payment"] = "Payment attributes are invalid"
	}

	return problems
}

func (d DeliveryAttrs) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if d.Name == "" {
		problems["name"] = "Name is required"
	}
	if d.Phone == "" {
		problems["phone"] = "Phone is required"
	}
	if d.Zip == "" {
		problems["zip"] = "Zip is required"
	}
	if d.City == "" {
		problems["city"] = "City is required"
	}
	if d.Address == "" {
		problems["address"] = "Address is required"
	}
	if d.Region == "" {
		problems["region"] = "Region is required"
	}
	if d.Email == "" {
		problems["email"] = "Email is required"
	}

	return problems
}

func (p PaymentAttrs) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if p.Transaction == "" {
		problems["transaction"] = "Transaction ID is required and must be a valid UUID"
	}
	if p.Currency == "" {
		problems["currency"] = "Currency is required"
	}
	if p.Provider == "" {
		problems["provider"] = "Provider is required"
	}
	if p.Amount <= 0 {
		problems["amount"] = "Amount must be a positive integer"
	}

	return problems
}
