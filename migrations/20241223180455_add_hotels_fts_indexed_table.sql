-- +goose Up
-- +goose StatementBegin
CREATE VIRTUAL TABLE IF NOT EXISTS hotels_fts USING fts5 (region_id, title);
INSERT INTO hotels_fts (region_id, title) SELECT region_id, title FROM hotels;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE hotels_fts;
-- +goose StatementEnd
