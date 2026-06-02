package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"hotels_api/internal/domain"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"
)

// SearchHotels method returns filtered request with hotels
func (s *SQLiteRepo) SearchHotels(ctx context.Context, text string, limit int) (*domain.Locations, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	result := &domain.Locations{
		Hotels: make(map[string][]string),
	}

	// full text search FTS5
	sqlStmt := sq.StatementBuilder.PlaceholderFormat(sq.Question)

	query, args, err := sqlStmt.Select("reg.title, fts.title").
		Distinct().
		From("hotels_fts AS fts").
		Join("regions AS reg ON fts.region_id = reg.region_id").
		Where("fts.title MATCH ?", text).
		// OrderBy("rank").
		Limit(uint64(limit)).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := s.DB.Query(query, args...)
	if err != nil || errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {

		var region, hotel string

		if err := rows.Scan(&region, &hotel); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		result.Hotels[region] = append(result.Hotels[region], hotel)
	}
	return result, nil
}
