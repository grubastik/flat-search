package db

import(
    "errors"
    "database/sql"
    "github.com/grubastik/flat-search/models"
)

var (
    // ErrAdvertDoesNotExist defines error for case advert is not exist in DB
    ErrAdvertDoesNotExist = errors.New("Advert does not exist")
)

// Advert implements db interface to store Advert
type Advert struct {}

// Insert inserts new record to db
func (a *Advert) Insert(e *models.Advert) error {
    var result sql.Result
    result, err := Storage.Run("INSERT INTO advert (locality, link, hash_id, price, name, description, status, created) VALUES( ?, ?, ?, ?, ?, ?, ?, ? )",
        e.Locality, e.Link, e.HashID, e.Price, e.Name, e.Description, e.Status, e.GetCreated().Unix())
    if err != nil {
        return err
    }
    e.ID, err = result.LastInsertId()
    return err
}

// Update updates existing record in DB
func (a *Advert) Update(e *models.Advert) error {
    _, err := Storage.Run("UPDATE advert SET locality = ?, link = ?, hash_id = ?, price = ?, name = ?, description = ?, status = ?, created = ? WHERE id = ?",
        e.Locality, e.Link, e.HashID, e.Price, e.Name, e.Description, e.Status, e.GetCreated().Unix(), e.ID)
    return err
}

// Delete removes existed record from db
func (a *Advert) Delete(id int64) error {
    _, err := Storage.Run("DELETE FROM advert WHERE id = ?", id)
    return err
}

// Load retrieves existed record from db
func (a *Advert) Load(f string, v interface{}, e *models.Advert) error {
    var row *sql.Row
    var stmt *sql.Stmt
    stmt, err := Storage.Prepare("SELECT id, locality, link, hash_id, price, name, description, status, created FROM advert WHERE " + f + " = ?")
    if err != nil {
        return err
    }
    defer stmt.Close()
    row = stmt.QueryRow(v)
    row.Scan(&e.ID, &e.Locality, &e.Link, &e.HashID, &e.Price, &e.Name, &e.Description, &e.Status, &e.Created)
    if e.ID == 0 {
        return ErrAdvertDoesNotExist
    }
    return nil
}
