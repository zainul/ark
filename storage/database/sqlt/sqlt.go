package sqlt

import (
	"database/sql"
	"log"
	"time"

	// MySQL
	_ "github.com/go-sql-driver/mysql"

	// PostgreSQL
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	tkpdSqlt "github.com/tokopedia/sqlt"
	"github.com/zainul/ark/storage/database"
)

// New sqlt module
func New(config Config) database.Database {

	// Setup connection string
	connectionString := config.Master

	// Get slave value if available
	if config.Slave != nil {
		for _, v := range config.Slave {
			connectionString += ";" + v
		}
	}

	// Open connection to DB
	db, err := tkpdSqlt.Open(config.Driver, connectionString)
	if err != nil {
		log.Println("func New", err)
		return nil
	}

	db.SetMaxOpenConnections(100)
	db.SetConnMaxLifetime(time.Minute * 10)

	return &sqltDB{
		config: config,
		db:     db,
	}
}

func (f *sqltDB) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	return f.db.Queryx(query, args...)
}

func (f *sqltDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return f.db.Exec(query, args...)
}
func (f *sqltDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return f.db.Query(query, args...)
}

func (f *sqltDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return f.db.QueryRow(query, args...)
}

func (f *sqltDB) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	return f.db.QueryRowx(query, args...)
}

func (f *sqltDB) Begin() (*sql.Tx, error) {
	return f.db.Begin()
}

func (f *sqltDB) Beginx() (*sqlx.Tx, error) {
	return f.db.Beginx()
}

func (f *sqltDB) Master() *sqlx.DB {
	return f.db.Master()
}

func (f *sqltDB) Get(dest interface{}, query string, args ...interface{}) error {
	return f.db.Get(dest, query, args...)
}

// SelectMaster using master db.
func (f *sqltDB) SelectMaster(dest interface{}, query string, args ...interface{}) error {
	return f.db.SelectMaster(dest, query, args...)
}
