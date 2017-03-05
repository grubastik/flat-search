package models

import(
    "database/sql"
    "./../db"
)

type Location struct {
    Id int64
    AdvertId int64
    Lat float64
    Lon float64
}

func (l *Location) SetId(id int64) {
    l.Id = id
}

func (l *Location) SetAdvertId(advertId int64) {
    l.AdvertId = advertId
}

func (l *Location) SetLatitude(lat float64) {
    l.Lat = lat
}

func (l *Location) SetLongitude(lon float64) {
    l.Lon = lon
}

func (l *Location) GetId() int64 {
    return l.Id
}

func (l *Location) GetAdvertId() int64 {
    return l.AdvertId
}

func (l *Location) GetLatitude() float64 {
    return l.Lat
}

func (l *Location) GetLongitude() float64 {
    return l.Lon
}

func (l *Location) Insert() {
    var result sql.Result
    var id int64
    result = db.Storage.Run("INSERT INTO locations (advert_id, lat, lon) VALUES( ?, ?, ? )",
        l.GetAdvertId(), l.GetLatitude(), l.GetLongitude())
    id, _ = result.LastInsertId()
    l.SetId(id)
}

func (l *Location) Update() {
    db.Storage.Run("UPDATE locations SET advert_id = ?, lat = ?, lon = ? WHERE id = ?",
        l.GetAdvertId(), l.GetLatitude(), l.GetLongitude(), l.GetId())
}

func (l *Location) Delete() {
    db.Storage.Run("DELETE FROM locations WHERE id = ?", l.GetId())
}

func (l *Location) Load() {
    var row *sql.Row
    var stmt *sql.Stmt
    stmt = db.Storage.Prepare("SELECT * FROM locations WHERE id = ?")
    defer stmt.Close()
    row = db.Storage.GetRow(stmt, l.GetId())
    row.Scan(l);
}

func (l *Location) LoadByAdvertId() {
    var row *sql.Row
    var stmt *sql.Stmt
    stmt = db.Storage.Prepare("SELECT * FROM locations WHERE advert_id = ?")
    defer stmt.Close()
    row = db.Storage.GetRow(stmt, l.GetAdvertId())
    row.Scan(l);
}
