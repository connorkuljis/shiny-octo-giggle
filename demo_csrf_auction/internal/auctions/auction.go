package auction

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Auction struct {
	Id          int
	Date        time.Time
	Title       string
	Description string
	Price       float64
}

func LoadMockAuctionData(db *sql.DB) []Auction {
	q := `SELECT * FROM auctions;`

	var records []Auction

	result, err := db.Query(q)
	if err != nil {
		panic(err)
	}

	for result.Next() {
		var auction Auction
		err := result.Scan(&auction.Id, &auction.Date, &auction.Title, &auction.Description, &auction.Price)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			break
		}
		records = append(records, auction)
	}

	return records
}

func GetAuctionById(db *sql.DB, id int) (Auction, error) {
	q := `SELECT * FROM auctions WHERE id = ?`

	result, err := db.Query(q, id)
	if err != nil {
		panic(err)
	}

	var auction Auction
	for result.Next() {
		err := result.Scan(&auction.Id, &auction.Date, &auction.Title, &auction.Description, &auction.Price)
		if err != nil {
			return auction, err
		}
	}

	return auction, nil

}

// Helper function to parse a date string into a time.Time object
func parseDate(dateStr string) time.Time {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		panic(err)
	}
	return date
}
