-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE withdrawals (
  id SERIAL,
  order_num text NOT NULL UNIQUE ,
  sum  numeric NOT NULL,
  processed_at timestamp default NOW(),
  user_id int,
  PRIMARY KEY(id)
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
