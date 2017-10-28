package db

import(
    "database/sql"
    "github.com/grubastik/flat-search/models"
)

// Property implements db interface to store properties linked to the advert
type Property struct {}

// Insert inserts new record to db
func (p *Property) Insert(e models.Property) error {
    var result sql.Result
    result, err := Storage.Run("INSERT INTO property (advert_id, name, value) VALUES( ?, ?, ? )", e.AdvertID, e.Name, e.Value)
    if err != nil {
        return err
    }
    e.ID, err = result.LastInsertId()
    return err
}

// Update updates existing record in DB
func (p *Property) Update(e *models.Property) error {
    _, err := Storage.Run("UPDATE property SET advert_id = ?, name = ?, value = ? WHERE id = ?", e.AdvertID, e.Name, e.Value, e.ID)
    return err
}

// Delete removes existed record from db
func (p *Property) Delete(id int64) error {
    _, err := Storage.Run("DELETE FROM property WHERE id = ?", id)
    return err
}

// Load retrieves existed records from db
func (p *Property) Load(f string, v interface{}, e *models.Property) ([]models.Property, error) {
    var (
        rows *sql.Rows
        stmt *sql.Stmt
        l []models.Property
    )
    
    stmt, err := Storage.Prepare("SELECT id, advert_id, name, value FROM property WHERE " + f + " = ?")
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
        err = rows.Scan(&ei.ID, &ei.AdvertID, &ei.Name, &ei.Value)
        if err != nil {
            return nil, err
        }
        if e.ID != 0 {
            l = append(l, ei)
        }
    }
    return l, nil
}
