package entity

import (
	"time"
)

type (
	Order struct {
		UID               string
		TrackNumber       string
		Entry             string
		Delivery          DeliveryAttrs
		Payment           PaymentAttrs
		Items             []ItemAttrs
		Locale            string
		InternalSignature string
		CustomerID        string
		DeliveryService   string
		ShardKey          string
		SmID              int
		DateCreated       time.Time
		OffShard          string
	}

	DeliveryAttrs struct {
		Name    string
		Phone   string
		Zip     string
		City    string
		Address string
		Region  string
		Email   string
	}
	PaymentAttrs struct {
		Transaction  string
		RequestID    string
		Currency     string
		Provider     string
		Amount       int
		PaymentDt    int
		Bank         string
		DeliveryCost int
		GoodsTotal   int
		CustomFee    int
	}
	ItemAttrs struct {
		ChrtID      int
		TrackNumber string
		Price       int
		Rid         string
		Name        string
		Sale        int
		Size        string
		TotalPrice  int
		NmID        int
		Brand       string
		Status      int
	}
)
