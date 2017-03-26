package models

import (
	"database/sql"
	"github.com/grubastik/flat-search/db"
)

// Location defines info about location of the advert
type Location struct {
	ID       int64
	AdvertID int64
	Lat      float64
	Lon      float64
}

// Insert inserts new record to db
func (l *Location) Insert() error {
	var result sql.Result
	var id int64
	result, err := db.Storage.Run("INSERT INTO location (advert_id, lat, lon) VALUES( ?, ?, ? )", l.AdvertID, l.Lat, l.Lon)
	if err != nil {
		return err
	}
	id, err = result.LastInsertId()
	if err != nil {
		return err
	}
	l.ID = id
	return nil
}

// Update updates existing record in DB
func (l *Location) Update() error {
	_, err := db.Storage.Run("UPDATE location SET advert_id = ?, lat = ?, lon = ? WHERE id = ?",
		l.AdvertID, l.Lat, l.Lon, l.ID)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes existed record from db
func (l *Location) Delete() error {
	_, err := db.Storage.Run("DELETE FROM location WHERE id = ?", l.ID)
	if err != nil {
		return err
	}
	return nil
}

// Load gets record from DB by id
func (l *Location) Load() error {
	return l.loadBy("id", l.ID)
}

// LoadByAdvertID  gets record from DB by AdvertID
func (l *Location) LoadByAdvertID() error {
	return l.loadBy("advert_id", l.AdvertID)
}

func (l *Location) loadBy(field string, value interface{}) error {
	var row *sql.Row
	var stmt *sql.Stmt
	stmt, err := db.Storage.Prepare("SELECT * FROM location WHERE " + field + " = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	row = stmt.QueryRow(value)
	row.Scan(l)
	return nil
}
