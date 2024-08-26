package pgdb

import (
	"context"

	"github.com/v7ktory/wb_task_one/internal/entity"
	"github.com/v7ktory/wb_task_one/pkg/postgres"
)

type Order interface {
	SaveOrder(ctx context.Context, order *entity.Order) (string, error)
	GetLRUOrders(ctx context.Context) ([]*entity.Order, error)
	UpdateOrderTime(ctx context.Context, uid string) error
}
type PgRepo struct {
	Order
}

func NewPgRepo(pg *postgres.Postgres) *PgRepo {
	return &PgRepo{
		Order: NewOrderRepo(pg),
	}
}
