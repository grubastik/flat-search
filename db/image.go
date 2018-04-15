package db

import (
	"database/sql"
	"github.com/grubastik/flat-search/models"
)

// Image implements db interface to store images linked to the advert
type Image struct{}

// Insert inserts new record to db
func (p *Image) Insert(e models.Image) error {
	var result sql.Result
	result, err := Storage.Run("INSERT INTO image (advert_id, url) VALUES( ?, ? )", e.AdvertID, e.URL)
	if err != nil {
		return err
	}
	e.ID, err = result.LastInsertId()
	return err
}

// Update updates existing record in DB
func (p *Image) Update(e *models.Image) error {
	_, err := Storage.Run("UPDATE image SET advert_id = ?, url = ? WHERE id = ?", e.AdvertID, e.URL, e.ID)
	return err
}

// Delete removes existed record from db
func (p *Image) Delete(id int64) error {
	_, err := Storage.Run("DELETE FROM image WHERE id = ?", id)
	return err
}

// Load retrieves existed records from db
func (p *Image) Load(f string, v interface{}, e *models.Image) ([]models.Image, error) {
	var (
		rows *sql.Rows
		stmt *sql.Stmt
		l    []models.Image
	)

	stmt, err := Storage.Prepare("SELECT id, advert_id, url FROM image WHERE " + f + " = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err = stmt.Query(v)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		ei := *e
		err = rows.Scan(&ei.ID, &ei.AdvertID, &ei.URL)
		if err != nil {
			return nil, err
		}
		if e.ID != 0 {
			l = append(l, ei)
		}
	}
	return l, nil
}
