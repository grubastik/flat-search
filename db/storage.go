package db

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "strconv"

    "github.com/grubastik/flat-search/config"
)

const configName = "db"

type storageEntity struct {
    db *sql.DB
}

var Storage *storageEntity

type Error interface {
	error
}

func NewDb(config *config.Config) (*storageEntity, error) {
    if (Storage != nil) {
        Storage.Close()
    }
    storage := new(storageEntity)
    err := storage.Connect(config.GetDb())
    if err != nil {
        return nil, err
    }
    return storage, nil
}

//in future make it configurable
func (st *storageEntity) Connect(config *config.Db) error {
    var err Error
    st.db, err = sql.Open(config.Engine, config.Username + ":" + config.Password + "@tcp(" + config.Host + ":" + strconv.Itoa(int(config.Port)) + ")/" + config.Database + "?charset=utf8")
    if err != nil {
        return nil
    }
    return nil
}

func (st *storageEntity) Close() {
    st.db.Close()
}

//accepts raw queries for insert update delete
func (st *storageEntity) Run(sql string, options... interface{}) (sql.Result, error) {
    return st.db.Exec(sql, options...)
}

//prepare stmt. Used for select and save
func (st *storageEntity) Prepare(sql string) (*sql.Stmt, error) {
    return st.db.Prepare(sql)
}

//exec stmt to get all rows
func (st *storageEntity) GetRows(stmt *sql.Stmt, options... interface{}) (*sql.Rows, error) {
    return stmt.Query(options...)
}
