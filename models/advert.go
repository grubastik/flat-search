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

var (
    ErrAdvertDoesNotExist = errors.New("Advert does not exist")
)

func NewAdvert() *Advert {
    var model = new(Advert);
    model.SetCreatedAt();
    model.SetStatusInitial();
    return model;
}

func (a *Advert) SetCreatedAt() {
    a.Created = time.Now().Unix();
}

func (a *Advert) SetStatusInitial() {
    a.Status = statusNew;
}

func (a *Advert) GetCreated() time.Time {
    return time.Unix(a.Created, 0)
}

func (a *Advert) ExistsInDbByHashId() (bool, error) {
    var row *sql.Row
    var stmt *sql.Stmt
    var test = new(Advert)
    stmt, err := db.Storage.Prepare("SELECT id FROM advert WHERE hash_id = ?")
    if err != nil {
        return false, err
    }
    defer stmt.Close()
    row = stmt.QueryRow(a.HashId)
    row.Scan(&test.Id)
    return test.Id > 0, nil
}

func (a *Advert) Insert() error {
    var result sql.Result
    var id int64
    result, err := db.Storage.Run("INSERT INTO advert (locality, link, hash_id, price, name, status, created) VALUES( ?, ?, ?, ?, ?, ?, ? )",
        a.Locality, a.Link, a.HashId, a.Price, a.Name, a.Status, a.GetCreated().Unix())
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

func (a *Advert) Update() error {
    _, err := db.Storage.Run("UPDATE advert SET locality = ?, link = ?, hash_id = ?, price = ?, name = ?, status = ? WHERE id = ?",
        a.Locality, a.Link, a.HashId, a.Price, a.Name, a.Status, a.GetCreated().Unix())
    if err != nil {
        return err
    }
    err = a.Location.Update()
    if err != nil {
        return err
    }
    return nil
}

func (a *Advert) Delete() error {
    err := a.Location.Delete()
    if err != nil {
        return err
    }
    _, err = db.Storage.Run("DELETE FROM advert WHERE id = ?", a.Id)
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
    stmt, err := db.Storage.Prepare("SELECT id, locality, link, hash_id, price, name, status, created FROM advert WHERE " + field + " = ?")
    if err != nil {
        return err
    }
    defer stmt.Close()
    row = stmt.QueryRow(value)
    row.Scan(&a.Id, &a.Locality, &a.Link, &a.HashId, &a.Price, &a.Name, &a.Status, &a.Created);
    if (a.Id > 0) {
        var lModel = new(Location);
        lModel.AdvertId = a.Id
        lModel.LoadByAdvertId()
        a.Location = lModel
    } else {
        return ErrAdvertDoesNotExist
    }
    return nil
}
