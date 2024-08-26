package pgdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/v7ktory/wb_task_one/internal/entity"
	"github.com/v7ktory/wb_task_one/pkg/postgres"
)

var ErrAlreadyExists = errors.New("order already exists")

type OrderRepo struct {
	*postgres.Postgres
}

func NewOrderRepo(pg *postgres.Postgres) *OrderRepo {
	return &OrderRepo{
		Postgres: pg,
	}
}

func (o *OrderRepo) Save(ctx context.Context, order *entity.Order) (string, error) {
	const op = "pgdb.order.go - Save"

	sql, args, _ := o.Builder.
		Insert("orders").
		Columns("order_uid,track_number,entry,delivery,payment,items,locale,internal_signature,customer_id,delivery_service,shardkey,sm_id,date_created,off_shard").
		Values(order.UID, order.TrackNumber, order.Entry, order.Delivery, order.Payment, order.Items, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OffShard).
		Suffix("RETURNING order_uid").
		ToSql()

	var uid string
	err := o.Pool.QueryRow(ctx, sql, args...).Scan(&uid)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return "", ErrAlreadyExists
			}
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return uid, nil
}
