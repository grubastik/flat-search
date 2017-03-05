package db

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "strconv"

    "./../error"
    "./../config"
)

const configName = "db"

type storageEntity struct {
    db *sql.DB
    stmt []*(sql.Stmt)
    rows []*(sql.Rows)
}

var Storage *storageEntity

func NewStorage(config *config.Config) (*storageEntity) {
    if (Storage != nil) {
        Storage.Close()
    }
    storage := new(storageEntity)
    moduleConfig := config.GetDb();
    storage.Connect(moduleConfig)
    return storage
}

//in future make it configurable
func (st *storageEntity) Connect(config *config.Db) {
    var conn *sql.DB
    conn, err := sql.Open(config.Engine, config.Username + ":" + config.Password + "@tcp(" + config.Host + ":" + strconv.Itoa(int(config.Port)) + ")/" + config.Database + "?charset=utf8")
    error.DebugError(err)
    st.db = conn
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
func (st *storageEntity) Run(sql string, options... interface{}) (sql.Result) {
    result, err := st.db.Exec(sql, options...)
    error.DebugError(err)
    return result
}

//prepare stmt. Used for select and save
func (st *storageEntity) Prepare(sql string) *sql.Stmt {
    stmt, err := st.db.Prepare(sql)
    error.DebugError(err)
    st.stmt = append(st.stmt, stmt)
    return stmt
}

//exec stmt for insert, create, delete
func (st *storageEntity) Execute(stmt *sql.Stmt, options... interface{}) (sql.Result) {
    result, err := stmt.Exec(options...)
    error.DebugError(err)
    return result
}

//accepts raw queries for insert update delete
func (st *storageEntity) GetData(sql string, options... interface{}) (*sql.Rows) {
    rows, err := st.db.Query(sql, options...)
    error.DebugError(err)
    st.rows = append(st.rows, rows)
    return rows
}

//exec stmt to get single row
func (st *storageEntity) GetRow(stmt *sql.Stmt, options... interface{}) (*sql.Row) {
    result := stmt.QueryRow(options...)
    return result
}

//exec stmt to get all rows
func (st *storageEntity) GetRows(stmt *sql.Stmt, options... interface{}) (*sql.Rows) {
    result, err := stmt.Query(options...)
    error.DebugError(err)
    st.rows = append(st.rows, result)
    return result
}
