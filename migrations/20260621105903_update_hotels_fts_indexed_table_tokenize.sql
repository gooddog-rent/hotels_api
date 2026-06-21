-- +goose Up
-- +goose StatementBegin

-- adding fts trigram index to table
CREATE VIRTUAL TABLE IF NOT EXISTS hotels_fts USING fts5(
    title,
    region_id UNINDEXED,
    tokenize="trigram remove_diacritics 1"
);

INSERT INTO hotels_fts(title, region_id)
SELECT title, region_id FROM hotels;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS hotels_fts;
-- +goose StatementEnd
