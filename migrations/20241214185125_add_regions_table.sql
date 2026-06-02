-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS regions (
    id integer UNIQUE PRIMARY KEY AUTOINCREMENT,
    region_id integer PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    created_at timestamp DEFAULT (CURRENT_TIMESTAMP),
    updated_at timestamp DEFAULT (CURRENT_TIMESTAMP)
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE regions;
-- +goose StatementEnd