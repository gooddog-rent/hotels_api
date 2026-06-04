package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gooddog-rent/hotels_api/internal/domain"

	"github.com/pressly/goose/v3"
)

// static interface implementation check for convinience
var _ Searcher = (*SQLiteRepo)(nil)

type Searcher interface {
	SearchHotels(ctx context.Context, text string, limit int) (*domain.Locations, error)
}

type Repositories struct {
	Searcher
	DBConnection *SQLiteRepo
}

// SQLiteRepo implementation
type SQLiteRepo struct {
	DB *sql.DB
	mu sync.RWMutex
}

func NewRepositories(HOTELS_PATH string) (*Repositories, error) {

	// init sqlite3 db
	dbRepo, err := NewSQLiteRepo(HOTELS_PATH)
	if err != nil {
		return nil, err
	}

	return &Repositories{
		Searcher:     dbRepo,
		DBConnection: dbRepo,
	}, nil
}

func NewSQLiteRepo(filepath string) (*SQLiteRepo, error) {

	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	const maxOpenConns = 20
	const maxIdleConns = 20
	const maxIdleTime = "15m"

	// Set the maximum number of open (in-use + idle) connections in the pool. Note that
	// passing a value less than or equal to 0 will mean there is no limit.
	db.SetMaxOpenConns(maxOpenConns)

	// Set the maximum number of idle connections in the pool. Again, passing a value
	// less than or equal to 0 will mean there is no limit.
	db.SetMaxIdleConns(maxIdleConns)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, fmt.Errorf("failed to parse maxIdleTime duration %w", err)
	}

	// Set the maximum idle timeout.
	db.SetConnMaxIdleTime(duration)

	const timeout = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	if err := InitDB(context.Background(), db); err != nil {
		log.Fatalln(err)
	}
	log.Println("Connected to the sqlite3 database")

	return &SQLiteRepo{
		DB: db,
	}, nil
}

func InitDB(ctx context.Context, db *sql.DB) error {

	err := goose.Up(db, "migrations")
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
