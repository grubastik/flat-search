package models

import (
	"database/sql"
	"errors"
	"github.com/grubastik/flat-search/db"
	"time"
)

const statusNew = "new"

// Advert defines which advert info will be stored/processed
type Advert struct {
	ID       int64
	Locality string
	Link     string
	HashID   int64
	Price    float64
	Name     string
	Status   string
	Created  int64
	Location *Location
}

var (
	// ErrAdvertDoesNotExist defines error for case advert is not exist in DB
	ErrAdvertDoesNotExist = errors.New("Advert does not exist")
)

// NewAdvert creates new structure for advert
func NewAdvert() *Advert {
	var model = new(Advert)
	model.SetCreatedAt()
	model.SetStatusInitial()
	return model
}

// SetCreatedAt set created to current time
func (a *Advert) SetCreatedAt() {
	a.Created = time.Now().Unix()
}

// SetStatusInitial set status to new
func (a *Advert) SetStatusInitial() {
	a.Status = statusNew
}

// GetCreated converts timestamp to time struct and returns it
func (a *Advert) GetCreated() time.Time {
	return time.Unix(a.Created, 0)
}

// ExistsInDbByHashID checks if advert exists in db based on HashID
func (a *Advert) ExistsInDbByHashID() (bool, error) {
	var test = new(Advert)
	stmt, err := db.Storage.Prepare("SELECT id FROM advert WHERE hash_id = ?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(a.HashID)
	row.Scan(&test.ID)
	return test.ID > 0, nil
}

// Insert inserts new record to db
func (a *Advert) Insert() error {
	var result sql.Result
	var id int64
	result, err := db.Storage.Run("INSERT INTO advert (locality, link, hash_id, price, name, status, created) VALUES( ?, ?, ?, ?, ?, ?, ? )",
		a.Locality, a.Link, a.HashID, a.Price, a.Name, a.Status, a.GetCreated().Unix())
	if err != nil {
		return err
	}
	id, err = result.LastInsertId()
	if err != nil {
		return err
	}
	a.ID = id
	a.Location.AdvertID = a.ID
	a.Location.Insert()
	return nil
}

// Update updates existing record in DB
func (a *Advert) Update() error {
	_, err := db.Storage.Run("UPDATE advert SET locality = ?, link = ?, hash_id = ?, price = ?, name = ?, status = ? WHERE id = ?",
		a.Locality, a.Link, a.HashID, a.Price, a.Name, a.Status, a.GetCreated().Unix())
	if err != nil {
		return err
	}
	err = a.Location.Update()
	if err != nil {
		return err
	}
	return nil
}

// Delete removes existed record from db
func (a *Advert) Delete() error {
	err := a.Location.Delete()
	if err != nil {
		return err
	}
	_, err = db.Storage.Run("DELETE FROM advert WHERE id = ?", a.ID)
	if err != nil {
		return err
	}
	return nil
}

// Load gets record from DB by id
func (a *Advert) Load() error {
	return a.loadBy("id", a.ID)
}

// LoadByHashID gets record from DB by HashID
func (a *Advert) LoadByHashID() error {
	return a.loadBy("hash_id", a.HashID)
}

func (a *Advert) loadBy(field string, value interface{}) error {
	var row *sql.Row
	var stmt *sql.Stmt
	stmt, err := db.Storage.Prepare("SELECT id, locality, link, hash_id, price, name, status, created FROM advert WHERE " + field + " = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	row = stmt.QueryRow(value)
	row.Scan(&a.ID, &a.Locality, &a.Link, &a.HashID, &a.Price, &a.Name, &a.Status, &a.Created)
	if a.ID > 0 {
		var lModel = new(Location)
		lModel.AdvertID = a.ID
		lModel.LoadByAdvertID()
		a.Location = lModel
	} else {
		return ErrAdvertDoesNotExist
	}
	return nil
}
