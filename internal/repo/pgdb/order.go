package pgdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
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

func (o *OrderRepo) SaveOrder(ctx context.Context, order *entity.Order) (string, error) {
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
		return "", fmt.Errorf("%s - Pool.QueryRow: %w", op, err)
	}

	return uid, nil
}

func (o *OrderRepo) GetLRUOrders(ctx context.Context) ([]*entity.Order, error) {
	const op = "pgdb.order.go - GetLRUOrders"

	sql, args, _ := o.Builder.
		Select("order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, off_shard").
		From("orders").
		OrderBy("created_at DESC").
		Limit(1_073_741_824).
		ToSql()

	var orders []*entity.Order
	rows, err := o.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s - Pool.Query: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		order := new(entity.Order)
		err := rows.Scan(
			&order.UID,
			&order.TrackNumber,
			&order.Entry,
			&order.Delivery,
			&order.Payment,
			&order.Items,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SmID,
			&order.DateCreated,
			&order.OffShard,
		)
		if err != nil {
			return nil, fmt.Errorf("%s - rows.Scan: %w", op, err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (o *OrderRepo) UpdateOrderTime(ctx context.Context, uid string) error {
	const op = "pgdb.order.go - UpdateOrderTime"

	sql, args, _ := o.Builder.
		Update("orders").
		Set("created_at", squirrel.Expr("now() AT TIME ZONE 'utc'")).
		Where("order_uid = ?", uid).
		ToSql()

	_, err := o.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s - Pool.Query: %w", op, err)
	}
	return nil
}
