-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE withdrawal (
  id SERIAL,
  order text NOT NULL,
  sum  numeric NOT NULL,
  processed_at timestamp default NOW(),
  user_id int,
  PRIMARY KEY(id)
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
