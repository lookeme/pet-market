-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE balance (
    current  numeric NOT NULL default 0,
    withdrawn numeric NOT NULL default 0,
    user_id int
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
