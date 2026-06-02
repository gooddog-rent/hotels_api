package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Your JSON file path and database path
const JSON_FILE_PATH = "./hotels/updated_hotels_booking.com.json"
const DATABASE_PATH = "./hotels_test_db.sql"

// Locations list
type Locations struct {
	Locations []Location `json:"locations"`
}

// Location to store regions and hotels list
type Location struct {
	Region string   `json:"region"`
	Hotels []string `json:"hotels"`
}

func main() {

	// open new connection database pool
	db, err := sql.Open("sqlite3", DATABASE_PATH)
	if err != nil {
		log.Fatalf("Cannot to open a database: %v", err)
	}
	defer db.Close()

	// create new hotels and regions tables
	createTableQuery := `DROP TABLE regions; 
						 DROP TABLE hotels;

						CREATE TABLE IF NOT EXISTS regions (
    						id integer UNIQUE PRIMARY KEY AUTOINCREMENT,
    						region_id integer NOT NULL,
    						title TEXT NOT NULL,
    						created_at timestamp DEFAULT (CURRENT_TIMESTAMP),
    						updated_at timestamp DEFAULT (CURRENT_TIMESTAMP)
						);
						
						CREATE TABLE IF NOT EXISTS hotels (
   							id integer UNIQUE PRIMARY KEY AUTOINCREMENT,
    						region_id integer NOT NULL,
    						title TEXT NOT NULL,
    						created_at timestamp DEFAULT (CURRENT_TIMESTAMP),
    						updated_at timestamp DEFAULT (CURRENT_TIMESTAMP),
							FOREIGN KEY (region_id) REFERENCES regions (region_id)
						);
						`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Error to create tables: %v", err)
	}

	// read json file
	jsonFile, err := os.Open(JSON_FILE_PATH)
	if err != nil {
		log.Fatalf("Cannot to open a json file: %v", err)
	}
	defer jsonFile.Close()

	var locations Locations

	err = json.NewDecoder(jsonFile).Decode(&locations)
	if err != nil {
		log.Fatalf("Error to decode json: %v", err)
	}

	ctx := context.Background()

	// transaction start
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Transaction begin failed: %v", err)
		return
	}

	committed := false

	defer func() {
		if !committed {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				slog.Error("failed to rollback transaction", "error", rollbackErr)
			}
		}
	}()

	stmtRegions, err := tx.Prepare(`INSERT INTO regions(region_id, title) VALUES (?, ?);`)
	if err != nil {
		log.Printf("Query prepare error: %v", err)
		return
	}
	defer stmtRegions.Close()

	stmtHotels, err := tx.Prepare(`INSERT INTO hotels (title, region_id) VALUES (?, ?);`)
	if err != nil {
		log.Printf("Query prepare error: %v", err)
		return
	}
	defer stmtHotels.Close()

	for index, location := range locations.Locations {

		region_id := index

		_, err = stmtRegions.Exec(region_id, location.Region)
		if err != nil {
			log.Printf("insert data error: %v\n", err)
			return
		}

		for _, hotel := range location.Hotels {

			_, err = stmtHotels.Exec(hotel, region_id)
			if err != nil {
				log.Printf("insert data error: %v\n", err)
				return
			}
		}
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Cannot to finished transaction: %v", err)
		return
	}
	committed = true

	fmt.Println("Successfully migrate from JSON file into SQLite3 database.")
}
