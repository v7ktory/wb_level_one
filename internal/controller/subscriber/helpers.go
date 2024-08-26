package subscriber

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/v7ktory/wb_task_one/internal/entity"
	"github.com/v7ktory/wb_task_one/internal/model"
)

func DecodeNATSReq[T model.Validator](data []byte) (T, map[string]string, error) {
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return v, nil, fmt.Errorf("decode json: %w", err)
	}
	if problems := v.Valid(context.Background()); len(problems) > 0 {
		return v, problems, fmt.Errorf("invalid %T: %d problems", v, len(problems))
	}
	return v, nil, nil
}

func ConvertNATSReq(orderRequest model.Order) *entity.Order {
	delivery := entity.DeliveryAttrs{
		Name:    orderRequest.Delivery.Name,
		Phone:   orderRequest.Delivery.Phone,
		Zip:     orderRequest.Delivery.Zip,
		City:    orderRequest.Delivery.City,
		Address: orderRequest.Delivery.Address,
		Region:  orderRequest.Delivery.Region,
		Email:   orderRequest.Delivery.Email,
	}
	payment := entity.PaymentAttrs{
		Transaction:  orderRequest.Payment.Transaction,
		RequestID:    orderRequest.Payment.RequestID,
		Currency:     orderRequest.Payment.Currency,
		Provider:     orderRequest.Payment.Provider,
		Amount:       orderRequest.Payment.Amount,
		PaymentDt:    orderRequest.Payment.PaymentDt,
		Bank:         orderRequest.Payment.Bank,
		DeliveryCost: orderRequest.Payment.DeliveryCost,
		GoodsTotal:   orderRequest.Payment.GoodsTotal,
		CustomFee:    orderRequest.Payment.CustomFee,
	}
	items := make([]entity.ItemAttrs, len(orderRequest.Items))
	for i, item := range orderRequest.Items {
		items[i] = entity.ItemAttrs{
			ChrtID:      item.ChrtID,
			TrackNumber: item.TrackNumber,
			Price:       item.Price,
			Rid:         item.Rid,
			Name:        item.Name,
			Sale:        item.Sale,
			Size:        item.Size,
			TotalPrice:  item.TotalPrice,
			NmID:        item.NmID,
			Brand:       item.Brand,
			Status:      item.Status,
		}
	}
	return &entity.Order{
		UID:               orderRequest.UID,
		TrackNumber:       orderRequest.TrackNumber,
		Entry:             orderRequest.Entry,
		Delivery:          delivery,
		Payment:           payment,
		Items:             items,
		Locale:            orderRequest.Locale,
		InternalSignature: orderRequest.InternalSignature,
		CustomerID:        orderRequest.CustomerID,
		DeliveryService:   orderRequest.DeliveryService,
		ShardKey:          orderRequest.ShardKey,
		SmID:              orderRequest.SmID,
		DateCreated:       orderRequest.DateCreated,
		OffShard:          orderRequest.OffShard,
	}
}
