package dummydb

import (
	"database/sql"

	// MySQL
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/tokopedia/sqlt"
	"github.com/zainul/ark/storage/database"
)

// New dummy database
func New(paramDb *sql.DB) database.Database {

	dbMock := sqlt.InitMocking(paramDb, 1)
	return &dummyDB{
		db: dbMock,
	}
}

func (f *dummyDB) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	return f.db.Queryx(query, args...)
}

func (f *dummyDB) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	return f.db.QueryRowx(query, args...)
}

func (f *dummyDB) Begin() (*sql.Tx, error) {
	return f.db.Begin()
}

func (f *dummyDB) Beginx() (*sqlx.Tx, error) {
	return f.db.Beginx()
}

func (f *dummyDB) Master() *sqlx.DB {
	return f.db.Master()
}

func (f *dummyDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return f.db.Exec(query, args...)
}

func (f *dummyDB) Get(dest interface{}, query string, args ...interface{}) error {
	return f.db.Get(dest, query, args...)
}

func (f *dummyDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return f.db.Query(query, args...)
}

func (f *dummyDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return f.db.QueryRow(query, args...)
}

// SelectMaster using master db.
func (f *dummyDB) SelectMaster(dest interface{}, query string, args ...interface{}) error {
	return f.db.SelectMaster(dest, query, args...)
}
