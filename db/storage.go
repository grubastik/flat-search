package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" //this package responsible for the sql queries - so we import it here
	"strconv"

	"github.com/grubastik/flat-search/config"
)

const configName = "db"

// StorageEntity contains connection dto DB
type StorageEntity struct {
	Db *sql.DB
}

// Storage contains Database connection
var Storage *StorageEntity

// NewDb creates new connection to the database
func NewDb(config *config.Config) (*StorageEntity, error) {
	if Storage != nil {
		Storage.Close()
	}
	Storage = new(StorageEntity)
	err := Storage.Connect(config.GetDb())
	return Storage, err
}

// Connect performs connection to the DB and store it locally
func (st *StorageEntity) Connect(config *config.Db) error {
	var err error
	st.Db, err = sql.Open(config.Engine, config.Username+":"+config.Password+"@tcp("+config.Host+":"+strconv.Itoa(int(config.Port))+")/"+config.Database+"?charset=utf8")
	return err
}

// Close closes connection to the DB
func (st *StorageEntity) Close() {
	st.Db.Close()
}

// Run accepts raw queries for insert update delete to run on opened connection
func (st *StorageEntity) Run(sql string, options ...interface{}) (sql.Result, error) {
	return st.Db.Exec(sql, options...)
}

// Prepare prepares statement. Used for select and save
func (st *StorageEntity) Prepare(sql string) (*sql.Stmt, error) {
	return st.Db.Prepare(sql)
}

// GetRows execs stmt to get all rows from DB based on prepared query
func (st *StorageEntity) GetRows(stmt *sql.Stmt, options ...interface{}) (*sql.Rows, error) {
	return stmt.Query(options...)
}
