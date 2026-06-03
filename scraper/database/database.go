package database

import (
	"database/sql"
	_ "database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB *sql.DB
}

func NewDB() *Database {
	db, err := sql.Open("sqlite3", "ram_prices.db")
	if err != nil {
		log.Fatal(err)
	}

	return &Database{DB: db}
}

// migrates the database and initializes shit
func (d *Database) Migrate() {
	queries := []string{

		`CREATE TABLE IF NOT EXISTS ram_prices(
    	ID INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE ,
    	MARKET TEXT NOT NULL ,
	    PRODUCT_ID INTEGER UNIQUE NOT NULL ,
	    RAM_TYPE TEXT NOT NULL ,
	    RAM_CAPACITY TEXT NOT NULL ,
	    PRICE INTEGER NOT NULL ,
	    REGION TEXT NOT NULL ,
	    BRAND TEXT NOT NULL,
	    URL TEXT NOT NULL UNIQUE ,
	    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`,
	}

	for _, q := range queries {
		_, err := d.DB.Exec(q)
		if err != nil {
			// ⚠️ SQLite will error if column already exists
			// so we just ignore those safely
			log.Println("migration warning:", err)
		}
	}
}

func (d *Database) Insert(ID int, market string, ramType string, ramCapacity string, price int, region string, brand string, url string) error {
	query := `
	INSERT INTO ram_prices(
    PRODUCT_ID,
    MARKET,
    RAM_TYPE,
    RAM_CAPACITY,
    PRICE,
    REGION,
    BRAND,
    URL
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(PRODUCT_ID) DO UPDATE SET
    PRICE = excluded.PRICE,
    REGION = excluded.REGION,
    BRAND = excluded.BRAND,
    URL = excluded.URL,
    RAM_TYPE = excluded.RAM_TYPE,
    RAM_CAPACITY = excluded.RAM_CAPACITY,
    MARKET = excluded.MARKET;`

	_, err := d.DB.Exec(query, ID, market, ramType, ramCapacity, price, region, brand, url)
	if err != nil {
		return err
	}

	return nil
}
