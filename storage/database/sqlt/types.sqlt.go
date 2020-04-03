package sqlt

import (
	tkpdSqlt "github.com/tokopedia/sqlt"
)

// Config database configuration
type Config struct {
	Driver string
	Master string
	Slave  []string
}

// database module
type sqltDB struct {
	config Config
	db     *tkpdSqlt.DB
}
