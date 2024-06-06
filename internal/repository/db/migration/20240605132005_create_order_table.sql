-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE order(
  order text,
  status text NOT NULL default "PROCESSING",
  accrual numeric,
  uploaded_at TIMESTAMP DEFAULT NOW(),
  user_id int NOT NULL,
  PRIMARY KEY(order)
);


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
