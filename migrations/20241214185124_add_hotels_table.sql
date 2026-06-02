-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS hotels (
    id integer UNIQUE PRIMARY KEY AUTOINCREMENT,
    region_id integer NOT NULL,
    title TEXT NOT NULL,
    created_at timestamp DEFAULT (CURRENT_TIMESTAMP),
    updated_at timestamp DEFAULT (CURRENT_TIMESTAMP),
    FOREIGN KEY (region_id) REFERENCES regions (region_id)
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE hotels;
-- +goose StatementEnd