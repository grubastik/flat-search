package models

import(
    "database/sql"
    "github.com/grubastik/flat-search/db"
)

type Location struct {
    Id int64
    AdvertId int64
    Lat float64
    Lon float64
}

func (l *Location) Insert() error {
    var result sql.Result
    var id int64
    result, err := db.Storage.Run("INSERT INTO location (advert_id, lat, lon) VALUES( ?, ?, ? )", l.AdvertId, l.Lat, l.Lon)
    if err != nil {
        return err
    }
    id, err = result.LastInsertId()
    if err != nil {
        return err
    }
    l.Id = id
    return nil
}

func (l *Location) Update() error {
    _, err := db.Storage.Run("UPDATE location SET advert_id = ?, lat = ?, lon = ? WHERE id = ?",
        l.AdvertId, l.Lat, l.Lon, l.Id)
    if err != nil {
        return err
    }
    return nil
}

func (l *Location) Delete() error {
    _, err := db.Storage.Run("DELETE FROM location WHERE id = ?", l.Id)
    if err != nil {
        return err
    }
    return nil
}

func (l *Location) Load() error {
    return l.loadBy("id", l.Id)
}

func (l *Location) LoadByAdvertId() error {
    return l.loadBy("advert_id", l.AdvertId)
}

func (l *Location) loadBy(field string, value interface{}) error {
    var row *sql.Row
    var stmt *sql.Stmt
    stmt, err := db.Storage.Prepare("SELECT * FROM location WHERE " + field + " = ?")
    if err != nil {
        return err
    }
    defer stmt.Close()
    row = stmt.QueryRow(value)
    row.Scan(l);
    return nil
}
