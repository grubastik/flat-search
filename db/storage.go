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
    stmt []*(sql.Stmt)
    rows []*(sql.Rows)
}

var Storage *storageEntity

func NewDb(config *config.Config) (*storageEntity, error) {
    if (Storage != nil) {
        Storage.Close()
    }
    storage := new(storageEntity)
    moduleConfig := config.GetDb();
    err := storage.Connect(moduleConfig)
    if err != nil {
        return nil, err
    }
    return storage, nil
}

//in future make it configurable
func (st *storageEntity) Connect(config *config.Db) (error) {
    var conn *sql.DB
    conn, err := sql.Open(config.Engine, config.Username + ":" + config.Password + "@tcp(" + config.Host + ":" + strconv.Itoa(int(config.Port)) + ")/" + config.Database + "?charset=utf8")
    if err != nil {
        return nil
    }
    st.db = conn
    return nil
}

func (st *storageEntity) Close() {
    st.db.Close()
    for _, rows := range(st.rows) {
        //ignore errors. probably already closed
        rows.Close()
    }
    for _, stmt := range(st.stmt) {
        //ignore errors. probably already closed
        stmt.Close()
    }
}

//accepts raw queries for insert update delete
func (st *storageEntity) Run(sql string, options... interface{}) (sql.Result, error) {
    result, err := st.db.Exec(sql, options...)
    if err != nil {
        return nil, err
    }
    return result, nil
}

//prepare stmt. Used for select and save
func (st *storageEntity) Prepare(sql string) (*sql.Stmt, error) {
    stmt, err := st.db.Prepare(sql)
    if err != nil {
        return nil, err
    }
    st.stmt = append(st.stmt, stmt)
    return stmt, nil
}

//exec stmt for insert, create, delete
func (st *storageEntity) Execute(stmt *sql.Stmt, options... interface{}) (sql.Result, error) {
    result, err := stmt.Exec(options...)
    if err != nil {
        return nil, err
    }
    return result, nil
}

//accepts raw queries for insert update delete
func (st *storageEntity) GetData(sql string, options... interface{}) (*sql.Rows, error) {
    rows, err := st.db.Query(sql, options...)
    if err != nil {
        return nil, err
    }
    st.rows = append(st.rows, rows)
    return rows, nil
}

//exec stmt to get single row
func (st *storageEntity) GetRow(stmt *sql.Stmt, options... interface{}) (*sql.Row) {
    result := stmt.QueryRow(options...)
    return result
}

//exec stmt to get all rows
func (st *storageEntity) GetRows(stmt *sql.Stmt, options... interface{}) (*sql.Rows, error) {
    result, err := stmt.Query(options...)
    if err != nil {
        return nil, err
    }
    st.rows = append(st.rows, result)
    return result, nil
}
