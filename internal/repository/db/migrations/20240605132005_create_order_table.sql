-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE orders (
  order_id text NOT NULL,
  status text NOT NULL default 'PROCESSING',
  accrual numeric,
  uploaded_at TIMESTAMP DEFAULT NOW(),
  user_id int NOT NULL,
  PRIMARY KEY(order_id)
);


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
