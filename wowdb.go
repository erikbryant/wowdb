package wowdb

import (
	"database/sql"
	"fmt"
	"log"
)

var (
	db *sql.DB
)

// Item contains the properties of a single auction house item.
type Item struct {
	ID        int64
	Name      string
	SellPrice int64
	JSON      string
}

// Auction contains the properties of a single auction house auction.
type Auction struct {
	Auc           int64
	Item          int64
	Owner         string
	Bid           int64
	Buyout        int64
	Quantity      int64
	TimeLeft      string
	Rand          int64
	Seed          int64
	Context       int64
	HasBonusLists bool
	HasModifiers  bool
	PetBreedID    int64
	PetLevel      int64
	PetQualityID  int64
	PetSpeciesID  int64
	JSON          string
}

// Open opens a connection to the database.
func Open() {
	var err error

	db, err = sql.Open("mysql", "wow:wowpassword@tcp(127.0.0.1:3306)/wow")
	if err != nil {
		log.Fatal(err.Error())
	}
}

// Close closes the connection to the database.
func Close() {
	db.Close()
}

// SaveItem writes the given item to the database.
func SaveItem(item Item) {
	sqlString := "INSERT IGNORE INTO items ( id, name, sellPrice, json ) VALUES ( ?, ?, ?, ? )"

	_, err := db.Exec(sqlString, item.ID, item.Name, item.SellPrice, item.JSON)
	if err != nil {
		fmt.Println("SaveItem Exec:", err, item)
	}
}

// LookupItem looks for the given ID in the database and returns it if found.
func LookupItem(id int64) (Item, bool) {
	var item Item

	sqlString := "SELECT * FROM items WHERE id = " + fmt.Sprintf("%d", id) + " LIMIT 1"

	rows := db.QueryRow(sqlString)
	err := rows.Scan(&item.ID, &item.Name, &item.SellPrice, &item.JSON)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("LookupItem Scan:", err, item)
		}
		return item, false
	}

	return item, true
}

// SaveAuction writes the given auction to the database.
func SaveAuction(auction Auction) {
	var sqlString string
	var err error

	// If the row already exists we want to update it.
	if current, ok := LookupAuction(auction.Auc); ok {
		if auction != current {
			sqlString = "UPDATE auctions SET item = ?, owner = ?, bid = ?, buyout = ?, quantity = ?, timeLeft = ?, rand = ?, seed = ?, context = ?, hasBonusLists = ?, hasModifiers = ?, petBreedId = ?, petLevel = ?, petQualityId = ?, petSpeciesId = ?, json = ? WHERE auc = ?"
			_, err = db.Exec(sqlString, auction.Item, auction.Owner, auction.Bid, auction.Buyout, auction.Quantity, auction.TimeLeft, auction.Rand, auction.Seed, auction.Context, auction.HasBonusLists, auction.HasModifiers, auction.PetBreedID, auction.PetLevel, auction.PetQualityID, auction.PetSpeciesID, auction.JSON, auction.Auc)
			if err != nil {
				fmt.Println("SaveAuction Exec(UPDATE):", err, auction)
			}
		}
	} else {
		sqlString = "INSERT INTO auctions ( auc, item, owner, bid, buyout, quantity, timeLeft, rand, seed, context, hasBonusLists, hasModifiers, PetBreedId, petLevel, petQualityId, petSpeciesId, json ) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )"
		_, err = db.Exec(sqlString, auction.Auc, auction.Item, auction.Owner, auction.Bid, auction.Buyout, auction.Quantity, auction.TimeLeft, auction.Rand, auction.Seed, auction.Context, auction.HasBonusLists, auction.HasModifiers, auction.PetBreedID, auction.PetLevel, auction.PetQualityID, auction.PetSpeciesID, auction.JSON)
		if err != nil {
			fmt.Println("SaveAuction Exec(INSERT):", err, auction)
		}
	}

}

// LookupAuction looks for the given auction in the database and returns it if found.
func LookupAuction(auc int64) (Auction, bool) {
	var auction Auction
	var lastUpdated string

	sqlString := "SELECT * FROM auctions WHERE auc = " + fmt.Sprintf("%d", auc) + " LIMIT 1"

	rows := db.QueryRow(sqlString)
	err := rows.Scan(&auction.Auc, &auction.Item, &auction.Owner, &auction.Bid, &auction.Buyout, &auction.Quantity, &auction.TimeLeft, &auction.Rand, &auction.Seed, &auction.Context, &auction.HasBonusLists, &auction.HasModifiers, &auction.PetBreedID, &auction.PetLevel, &auction.PetQualityID, &auction.PetSpeciesID, &auction.JSON, &lastUpdated)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("LookupAuction Scan:", err)
		}
		return auction, false
	}

	return auction, true
}

// countRows returns the number of rows in the given table.
func countRows(table string) int64 {
	var count int64

	sqlString := "SELECT count(*) FROM " + table
	rows := db.QueryRow(sqlString)
	err := rows.Scan(&count)
	if err != nil {
		fmt.Println("countRows Scan:", err)
		return -1
	}

	return count
}

// CountItems returns the number of items stored in the database.
func CountItems() int64 {
	return countRows("items")
}

// CountAuctions returns the number of items stored in the database.
func CountAuctions() int64 {
	return countRows("auctions")
}
