-- +goose Up
-- +goose StatementBegin
CREATE TABLE "orders" (
  "order_uid" varchar(255) PRIMARY KEY,
  "track_number" varchar(255) NOT NULL,
  "entry" varchar(255) NOT NULL,
  "delivery" jsonb NOT NULL,
  "payment" jsonb NOT NULL,
  "items" jsonb NOT NULL,
  "locale" varchar(255) NOT NULL,
  "internal_signature" varchar(255) NOT NULL,
  "customer_id" varchar(255) NOT NULL,
  "delivery_service" varchar(255) NOT NULL,
  "shardkey" varchar(255) NOT NULL,
  "sm_id" int NOT NULL,
  "date_created" timestamp NOT NULL,
  "off_shard" varchar(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "orders";
-- +goose StatementEnd
