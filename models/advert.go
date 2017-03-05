package models

import (
    "time"
    "errors"
    "database/sql"
    "./../db"
)

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

func (a *Advert) SetId(id int64) {
    a.Id = id
}

func (a *Advert) SetLocality(locality string) {
    a.Locality = locality
}

func (a *Advert) SetLink(link string) {
    a.Link = link
}

func (a *Advert) SetHash(hash int64) {
    a.HashId = hash
}

func (a *Advert) SetPrice(price float64) {
    a.Price = price
}

func (a *Advert) SetName(name string) {
    a.Name = name
}

func (a *Advert) SetStatus(status string) {
    a.Status = status
}

func (a *Advert) SetLocation(location *Location) {
    a.Location = location
}

func (a *Advert) GetId() int64 {
    return a.Id
}

func (a *Advert) GetLocality() string {
    return a.Locality
}

func (a *Advert) GetLink() string {
    return a.Link
}

func (a *Advert) GetHash() int64 {
    return a.HashId
}

func (a *Advert) GetPrice() float64 {
    return a.Price
}

func (a *Advert) GetName() string {
    return a.Name
}

func (a *Advert) GetStatus() string {
    if (a.Status == "") {
        a.Status = "new"
    }
    return a.Status
}

func (a *Advert) GetCreated() time.Time {
    if (a.Created == 0) {
        a.Created = time.Now().Unix()
    }
    return time.Unix(a.Created, 0)
}

func (a *Advert) GetLocation() *Location {
    return a.Location
}

func (a *Advert) ExistsInDbByHashId() bool {
    var row *sql.Row
    var stmt *sql.Stmt
    var test = new(Advert)
    stmt = db.Storage.Prepare("SELECT id FROM adverts WHERE hash_id = ?")
    defer stmt.Close()
    row = db.Storage.GetRow(stmt, a.GetHash())
    row.Scan(&test.Id)
    return test.GetId() > 0
}

func (a *Advert) Insert() {
    var result sql.Result
    var id int64
    result = db.Storage.Run("INSERT INTO adverts (locality, link, hash_id, price, name, status, created) VALUES( ?, ?, ?, ?, ?, ?, ? )",
        a.GetLocality(), a.GetLink(), a.GetHash(), a.GetPrice(), a.GetName(), a.GetStatus(), a.GetCreated().Unix())
    id, _ = result.LastInsertId()
    a.SetId(id)
    a.GetLocation().SetAdvertId(a.GetId())
    a.GetLocation().Insert()
}

func (a *Advert) Update() {
    db.Storage.Run("UPDATE adverts SET locality = ?, link = ?, hash_id = ?, price = ?, name = ?, status = ? WHERE id = ?",
        a.GetLocality(), a.GetLink(), a.GetHash(), a.GetPrice(), a.GetName(), a.GetStatus(), a.GetCreated().Unix())
    a.GetLocation().Update()
}

func (a *Advert) Delete() {
    a.GetLocation().Delete()
    db.Storage.Run("DELETE FROM adverts WHERE id = ?", a.GetId())
}

func (a *Advert) Load() error {
    var row *sql.Row
    var stmt *sql.Stmt
    stmt = db.Storage.Prepare("SELECT id, locality, link, hash_id, price, name, status, created FROM adverts WHERE id = ?")
    defer stmt.Close()
    row = db.Storage.GetRow(stmt, a.GetId())
    row.Scan(&a.Id, &a.Locality, &a.Link, &a.HashId, &a.Price, &a.Name, &a.Status, &a.Created);
    if (a.GetId() > 0) {
        var lModel = new(Location);
        lModel.SetAdvertId(a.GetId())
        lModel.LoadByAdvertId()
        a.SetLocation(lModel)
    } else {
        return errors.New("Advert does not exist")
    }
    return nil
}

func (a *Advert) LoadByHashId() error {
    var row *sql.Row
    var stmt *sql.Stmt
    stmt = db.Storage.Prepare("SELECT id, locality, link, hash_id, price, name, status, created FROM adverts WHERE hash_id = ?")
    defer stmt.Close()
    row = db.Storage.GetRow(stmt, a.GetHash())
    row.Scan(&a.Id, &a.Locality, &a.Link, &a.HashId, &a.Price, &a.Name, &a.Status, &a.Created);
    if (a.GetId() > 0) {
        var lModel = new(Location);
        lModel.SetAdvertId(a.GetId())
        lModel.LoadByAdvertId()
        a.SetLocation(lModel)
    } else {
        return errors.New("Advert does not exist")
    }
    return nil
}
