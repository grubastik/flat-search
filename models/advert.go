package models

import (
    "time"
    "errors"
    "database/sql"
    "github.com/grubastik/flat-search/db"
)

const statusNew = "new";

type Advert struct {
    Id int64
    Locality string
    Link string
    HashId int64
    Price float64
    Name string
    Status string
    Created int64
    Location *Location
}

func (a *Advert) GetStatus() string {
    if (a.Status == "") {
        a.Status = statusNew
    }
    return a.Status
}

func (a *Advert) GetCreated() time.Time {
    if (a.Created == 0) {
        a.Created = time.Now().Unix()
    }
    return time.Unix(a.Created, 0)
}

func (a *Advert) ExistsInDbByHashId() (bool, error) {
    var row *sql.Row
    var stmt *sql.Stmt
    var test = new(Advert)
    stmt, err := db.Storage.Prepare("SELECT id FROM adverts WHERE hash_id = ?")
    if err != nil {
        return false, err
    }
    defer stmt.Close()
    row = db.Storage.GetRow(stmt, a.HashId)
    row.Scan(&test.Id)
    return test.Id > 0, nil
}

func (a *Advert) Insert() (error) {
    var result sql.Result
    var id int64
    result, err := db.Storage.Run("INSERT INTO adverts (locality, link, hash_id, price, name, status, created) VALUES( ?, ?, ?, ?, ?, ?, ? )",
        a.Locality, a.Link, a.HashId, a.Price, a.Name, a.GetStatus(), a.GetCreated().Unix())
    if err != nil {
        return err
    }
    id, err = result.LastInsertId()
    if err != nil {
        return err
    }
    a.Id = id
    a.Location.AdvertId = a.Id
    a.Location.Insert()
    return nil
}

func (a *Advert) Update() (error) {
    _, err := db.Storage.Run("UPDATE adverts SET locality = ?, link = ?, hash_id = ?, price = ?, name = ?, status = ? WHERE id = ?",
        a.Locality, a.Link, a.HashId, a.Price, a.Name, a.GetStatus(), a.GetCreated().Unix())
    if err != nil {
        return err
    }
    err = a.Location.Update()
    if err != nil {
        return err
    }
    return nil
}

func (a *Advert) Delete() (error) {
    err := a.Location.Delete()
    if err != nil {
        return err
    }
    _, err = db.Storage.Run("DELETE FROM adverts WHERE id = ?", a.Id)
    if err != nil {
        return err
    }
    return nil
}

func (a *Advert) Load() error {
    return a.loadBy("id", a.Id)
}

func (a *Advert) LoadByHashId() error {
    return a.loadBy("hash_id", a.HashId)
}

func (a *Advert) loadBy(field string, value interface{}) error {
    var row *sql.Row
    var stmt *sql.Stmt
    stmt, err := db.Storage.Prepare("SELECT id, locality, link, hash_id, price, name, status, created FROM adverts WHERE " + field + " = ?")
    if err != nil {
        return err
    }
    defer stmt.Close()
    row = db.Storage.GetRow(stmt, value)
    row.Scan(&a.Id, &a.Locality, &a.Link, &a.HashId, &a.Price, &a.Name, &a.Status, &a.Created);
    if (a.Id > 0) {
        var lModel = new(Location);
        lModel.AdvertId = a.Id
        lModel.LoadByAdvertId()
        a.Location = lModel
    } else {
        return errors.New("Advert does not exist")
    }
    return nil
}
