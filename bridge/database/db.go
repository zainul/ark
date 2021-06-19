package database

import (
	"context"
	"database/sql"
)

//DriverImplementor is chooser implementation of DB interface
type DriverImplementor string

//DBType is type definition of database type like master slave
type DBType string

//TracerImplementor is type of implementation of tracer
type TracerImplementor string

const (
	// GORM is constant package of jinzhu gorm
	GORM DriverImplementor = "gorm"
	// SQLX is constant package of sqlx jomoiron
	SQLX   DriverImplementor = "sqlx"
	slave  DBType            = "slave"
	master DBType            = "master"

	ElasticAPM TracerImplementor = "es-apm"
)

type (
	// Config is configuration of db value should be sql.DB
	Config struct {
		// MasterConn is master connection of database
		MasterConn *sql.DB
		// Slave is master connection of database
		SlaveConn *sql.DB
		//LogMode is mode for log enable
		LogMode bool

		Tracer TracerImplementor
	}
)

var (
	driversImplementors map[DriverImplementor]func(ctx context.Context, dialect string, cfg Config) DB
)

func init() {
	ds := make(map[DriverImplementor]func(ctx context.Context, dialect string, cfg Config) DB)
	ds[GORM] = newGorm
	driversImplementors = ds
}

//NewDB is set instance of database
// parameter is context, dialect, config og database, and implementor
func NewDB(ctx context.Context, dialect string, cfg Config, driverPlugin DriverImplementor) DB {
	return driversImplementors[driverPlugin](ctx, dialect, cfg)
}

//go:generate mockgen -destination mock/mock_db.go github.com/zainul/ark/bridge/database DB
type DB interface {
	//RunAsUnit is run sequentially transaction
	RunAsUnit(action func(tx interface{}) error) error
	// Create is definitely insert to one table
	Create(ctx context.Context, data interface{}, txs ...interface{}) error
	// Update is definitely update to one table, with condition parameter (operator is AND if multi condition)
	Update(ctx context.Context, table string, data map[string]interface{}, whereCondition map[string]interface{}, txs ...interface{}) error
	// Delete is definitely delete to one table, with condition parameter (operator is AND if multi condition)
	Delete(ctx context.Context, table string, whereCondition map[string]interface{}, txs ...interface{}) error
	// QueryExec is raw query to exec statement of insert, update or delete
	QueryExec(ctx context.Context, txs interface{}, query string, args ...interface{}) error
	// QueryRaw is raw query to exec statement of insert, update or delete
	QueryRaw(ctx context.Context, txs interface{}, target interface{}, sql string, values ...interface{}) error
	// EntityBy is find in entity by one field and targeted result set in parameter
	EntityBy(ctx context.Context, field string, value interface{}, target interface{}, txs ...interface{}) error
}
