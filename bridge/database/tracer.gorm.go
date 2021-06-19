package database

import (
	"context"
	"database/sql"

	"github.com/jinzhu/gorm"
	"go.elastic.co/apm/module/apmgorm"
)

func GetGormTracerOpen(tracer TracerImplementor, dialect string, conn *sql.DB) (*gorm.DB, error) {
	if tracer == ElasticAPM {
		return apmgorm.Open(dialect, conn)
	}
	return gorm.Open(dialect, conn)
}

func GetGormTracerWithContext(ctx context.Context, tracer TracerImplementor, db *gorm.DB) *gorm.DB {
	if tracer == ElasticAPM {
		return apmgorm.WithContext(ctx, db)
	}

	return db
}
