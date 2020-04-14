package database

import (
	"context"
	"database/sql"
)

//DriverImplementor is chooser implementation of DB interface
type DriverImplementor string
//DBType is type definition of database type like master slave
type DBType string

const (
	// GORM is constant package of jinzhu gorm
	GORM   DriverImplementor = "gorm"
	// SQLX is constant package of sqlx jomoiron
	SQLX   DriverImplementor = "sqlx"
	slave  DBType            = "slave"
	master DBType            = "master"
)

type (
	// Config is configuration of db value should be sql.DB
	Config struct {
		// MasterConn is master connection of database
		MasterConn *sql.DB
		// Slave is master connection of database
		SlaveConn  *sql.DB
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
	// Create is definitely insert to one table
	Create(ctx context.Context, data interface{}) error
	// Update is definitely update to one table, with condition parameter (operator is AND if multi condition)
	Update(ctx context.Context, table string, data map[string]interface{}, whereCondition map[string]interface{}) error
	// Delete is definitely delete to one table, with condition parameter (operator is AND if multi condition)
	Delete(ctx context.Context, table string, whereCondition map[string]interface{}) error
	// QueryExec is raw query to exec statement of insert, update or delete
	QueryExec(ctx context.Context, query string, args ...interface{}) error
	// QueryRaw is raw query to exec statement of insert, update or delete
	QueryRaw(ctx context.Context, target interface{}, sql string, values ...interface{}) error
	// EntityBy is find in entity by one field and targeted result set in parameter
	EntityBy(ctx context.Context, field string, value interface{}, target interface{}) error
}
