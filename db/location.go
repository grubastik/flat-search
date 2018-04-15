package db

import (
	"database/sql"
	"errors"
	"github.com/grubastik/flat-search/models"
)

var (
	// ErrLocationDoesNotExist defines error for case location is not exist in DB
	ErrLocationDoesNotExist = errors.New("Location does not exist")
)

// Location implements db interface to store Advert
type Location struct{}

// Insert inserts new record to db
func (l *Location) Insert(e *models.Location) error {
	var result sql.Result
	result, err := Storage.Run("INSERT INTO location (advert_id, lat, lon) VALUES( ?, ?, ? )", e.AdvertID, e.Lat, e.Lon)
	if err != nil {
		return err
	}
	e.ID, err = result.LastInsertId()
	return err
}

// Update updates existing record in DB
func (l *Location) Update(e *models.Location) error {
	_, err := Storage.Run("UPDATE location SET advert_id = ?, lat = ?, lon = ? WHERE id = ?", e.AdvertID, e.Lat, e.Lon, e.ID)
	return err
}

// Delete removes existed record from db
func (l *Location) Delete(id int64) error {
	_, err := Storage.Run("DELETE FROM location WHERE id = ?", id)
	return err
}

// Load retrieves existed record from db
func (l *Location) Load(f string, v interface{}, e *models.Location) error {
	var row *sql.Row
	var stmt *sql.Stmt
	stmt, err := Storage.Prepare("SELECT id, advert_id, lat, lon FROM location WHERE " + f + " = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	row = stmt.QueryRow(v)
	row.Scan(&e.ID, &e.AdvertID, &e.Lat, &e.Lon)
	if e.ID == 0 {
		return ErrLocationDoesNotExist
	}
	return nil
}
