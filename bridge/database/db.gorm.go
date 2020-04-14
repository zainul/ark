package database

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jinzhu/gorm"
	"go.elastic.co/apm/module/apmgorm"
)

type jinzhuGorm struct {
	dbs map[DBType]*gorm.DB
}

func newGorm(ctx context.Context, dialect string, cfg Config) DB {
	var err error
	dbs := make(map[DBType]*gorm.DB)
	dbs[master], err = apmgorm.Open(dialect, cfg.MasterConn)

	if err != nil {
		log.Fatal("Failed to init db", err)
	}

	dbs[slave], err = gorm.Open(dialect, cfg.SlaveConn)

	if err != nil {
		log.Fatal("Failed to init db", err)
	}

	dbs[master].LogMode(true)
	dbs[slave].LogMode(true)

	return &jinzhuGorm{
		dbs: dbs,
	}

}

func (g *jinzhuGorm) Create(ctx context.Context, data interface{}) error {
	g.dbs[master] = apmgorm.WithContext(ctx, g.dbs[master])
	return g.dbs[master].Create(data).Error
}

func (g *jinzhuGorm) EntityBy(ctx context.Context, field string, value interface{}, target interface{}) error {
	g.dbs[slave] = apmgorm.WithContext(ctx, g.dbs[slave])
	return g.dbs[slave].Where(map[string]interface{}{field: value}).Find(target).Error
}

func (g *jinzhuGorm) Update(ctx context.Context, table string, data map[string]interface{}, whereCondition map[string]interface{}) error {
	query := fmt.Sprintf("UPDATE %s SET ", table)

	vals := make([]interface{}, 0)
	queryField := make([]string, 0)
	whereVals := make([]interface{}, 0)
	whereField := make([]string, 0)

	for key, valmap := range data {
		vals = append(vals, valmap)
		queryField = append(queryField, key+"=?")
	}

	query = query + strings.Join(queryField, " ,") + " WHERE "

	for k, v := range whereCondition {
		whereVals = append(whereVals, v)
		whereField = append(whereField, k+"=?")
	}

	query = query + strings.Join(whereField, " AND ")

	vals = append(vals, whereVals...)

	g.dbs[master] = apmgorm.WithContext(ctx, g.dbs[master])
	return g.dbs[master].Exec(query, vals...).Error
}

func (g *jinzhuGorm) Delete(ctx context.Context, table string, whereCondition map[string]interface{}) error {
	query := fmt.Sprintf("DELETE %s WHERE", table)
	whereVals := make([]interface{}, 0)
	whereField := make([]string, 0)

	for k, v := range whereCondition {
		whereVals = append(whereVals, v)
		whereField = append(whereField, k+"=?")
	}

	query = query + strings.Join(whereField, " AND ")
	g.dbs[master] = apmgorm.WithContext(ctx, g.dbs[master])
	return g.dbs[master].Exec(query, whereVals...).Error
}

func (g *jinzhuGorm) QueryExec(ctx context.Context, query string, args ...interface{}) error {
	g.dbs[master] = apmgorm.WithContext(ctx, g.dbs[master])
	return g.dbs[master].Exec(query, args...).Error
}

func (g *jinzhuGorm) QueryRaw(ctx context.Context, target interface{}, sql string, values ...interface{}) error {
	g.dbs[master] = apmgorm.WithContext(ctx, g.dbs[master])
	return g.dbs[slave].Raw(sql, values).Find(target).Error
}
