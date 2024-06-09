-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE users (
    id SERIAL,
    login text NOT NULL,
    pass  text NOT NULL,
    date_create timestamp default NOW(),
    is_active bool default true,
    PRIMARY KEY(id)
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
