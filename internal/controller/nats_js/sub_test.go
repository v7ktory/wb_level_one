package natsjs

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/v7ktory/wb_task_one/internal/controller/mocks"
	"github.com/v7ktory/wb_task_one/internal/entity"
)

func TestHandleMessage(t *testing.T) {
	mockOrder := mocks.NewOrder(t)
	mockCache := mocks.NewCache[string, *entity.Order](t)
	mockLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	validJSON := `{
                    "order_uid": "b563feb7b2b84b6test",
                    "track_number": "WBILMTESTTRACK",
                    "entry": "WBIL",
                    "delivery": {
                        "name": "Test Testov",
                        "phone": "+9720000000",
                        "zip": "2639809",
                        "city": "Kiryat Mozkin",
                        "address": "Ploshad Mira 15",
                        "region": "Kraiot",
                        "email": "test@gmail.com"
                    },
                    "payment": {
                        "transaction": "b563feb7b2b84b6test",
                        "request_id": "",
                        "currency": "USD",
                        "provider": "wbpay",
                        "amount": 1817,
                        "payment_dt": 1637907727,
                        "bank": "alpha",
                        "delivery_cost": 1500,
                        "goods_total": 317,
                        "custom_fee": 0
                    },
                    "items": [
                        {
                            "chrt_id": 9934930,
                            "track_number": "WBILMTESTTRACK",
                            "price": 453,
                            "rid": "ab4219087a764ae0btest",
                            "name": "Mascaras",
                            "sale": 30,
                            "size": "0",
                            "total_price": 317,
                            "nm_id": 2389212,
                            "brand": "Vivienne Sabo",
                            "status": 202
                        }
                    ],
                    "locale": "en",
                    "internal_signature": "",
                    "customer_id": "test",
                    "delivery_service": "meest",
                    "shardkey": "9",
                    "sm_id": 99,
                    "date_created": "2021-11-26T06:22:19Z",
                    "oof_shard": "1"
                }`

	type args struct {
		ctx context.Context
		msg []byte
	}

	tests := []struct {
		name      string
		args      args
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Test valid message",
			args: args{
				ctx: context.Background(),
				msg: []byte(validJSON),
			},
			mockSetup: func() {
				mockOrder.
					On("SaveOrder", mock.Anything, mock.AnythingOfType("*entity.Order")).
					Return("b563feb7b2b84b6test", nil)

				mockCache.
					On("Put", mock.Anything, mock.Anything).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Test invalid message",
			args: args{
				ctx: context.Background(),
				msg: []byte(`Invalid message`),
			},
			mockSetup: func() {
				// No mock setup needed for invalid message
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s := Subscriber{
			jetStr:    nil,
			orderRepo: mockOrder,
			cache:     mockCache,
			logger:    mockLogger,
		}
		tt.mockSetup()
		err := s.handleMessage(tt.args.ctx, tt.args.msg)
		if (err != nil) != tt.wantErr {
			t.Errorf("handleMessage() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
	}
}
