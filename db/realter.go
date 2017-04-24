package db

import(
    "errors"
    "database/sql"
    "github.com/grubastik/flat-search/models"
)

var (
    // ErrRealtorDoesNotExist defines error for case advert is not exist in DB
    ErrRealtorDoesNotExist = errors.New("Realtor does not exist")
)


// Realtor implements db interface to store realtor linked to the advert
type Realtor struct {}

// Insert inserts new record to db
func (l *Realtor) Insert(e models.Realtor) error {
    var result sql.Result
    result, err := Storage.Run("INSERT INTO realtor (advert_id, name, phone, email, company, company_phone, company_ico) VALUES( ?, ?, ?, ?, ?, ?, ? )", e.AdvertID, e.Name, e.Phone, e.Email, e.Company, e.CompanyPhone, e.CompanyICO)
    if err != nil {
        return err
    }
    e.ID, err = result.LastInsertId()
    return err
}

// Update updates existing record in DB
func (l *Realtor) Update(e *models.Realtor) error {
    _, err := Storage.Run("UPDATE realtor SET advert_id = ?, name = ?, phone = ?, email = ?, company = ?, company_phone = ?, company_ico = ? WHERE id = ?", e.AdvertID, e.Name, e.Phone, e.Email, e.Company, e.CompanyPhone, e.CompanyICO, e.ID)
    return err
}

// Delete removes existed record from db
func (l *Realtor) Delete(id int64) error {
    _, err := Storage.Run("DELETE FROM realtor WHERE id = ?", id)
    return err
}

// Load retrieves existed record from db
func (l *Realtor) Load(f string, v interface{}, e *models.Realtor) error {
    var row *sql.Row
    var stmt *sql.Stmt
    stmt, err := Storage.Prepare("SELECT id, advert_id, name, phone, email, company, company_phone, company_ico FROM location WHERE " + f + " = ?")
    if err != nil {
        return err
    }
    defer stmt.Close()
    row = stmt.QueryRow(v)
    row.Scan(&e.ID, &e.AdvertID, &e.Name, &e.Phone, &e.Email, &e.Company, &e.CompanyPhone, &e.CompanyICO)
    if e.ID == 0 {
        return ErrRealtorDoesNotExist
    }
    return nil
}
