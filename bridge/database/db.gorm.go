package database

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jinzhu/gorm"
)

type jinzhuGorm struct {
	dbMaster   *gorm.DB
	dbFollower *gorm.DB
	tracer     TracerImplementor
}

func newGorm(ctx context.Context, dialect string, cfg Config) DB {
	var err error
	var dbM, dbS *gorm.DB

	dbM, err = GetGormTracerOpen(cfg.Tracer, dialect, cfg.MasterConn)

	if err != nil {
		log.Fatal("Failed to init db", err)
	}

	dbS, err = GetGormTracerOpen(cfg.Tracer, dialect, cfg.MasterConn)

	if err != nil {
		log.Fatal("Failed to init db", err)
	}

	dbM.LogMode(cfg.LogMode)
	dbS.LogMode(cfg.LogMode)

	return &jinzhuGorm{
		dbMaster:   dbM,
		dbFollower: dbS,
		tracer:     cfg.Tracer,
	}

}

func (g *jinzhuGorm) RunAsUnit(action func(tx interface{}) error) error {
	var errAction error
	errAction = g.dbMaster.Transaction(func(tx *gorm.DB) error {
		if errAction = action(tx); errAction != nil {
			return errAction
		}
		return nil
	})
	return errAction
}

func (g *jinzhuGorm) Create(ctx context.Context, data interface{}, txs ...interface{}) error {

	g.dbMaster = GetGormTracerWithContext(ctx, g.tracer, g.dbMaster)

	if len(txs) > 0 && txs[0] != nil {
		tx, _ := txs[0].(*gorm.DB)
		return tx.Create(data).Error
	}
	return g.dbMaster.Create(data).Error
}

func (g *jinzhuGorm) EntityBy(ctx context.Context, field string, value interface{}, target interface{}, txs ...interface{}) error {

	g.dbFollower = GetGormTracerWithContext(ctx, g.tracer, g.dbMaster)

	if len(txs) > 0 && txs[0] != nil {
		tx, _ := txs[0].(*gorm.DB)
		return tx.Where(map[string]interface{}{field: value}).Find(target).Error
	}

	return g.dbFollower.Where(map[string]interface{}{field: value}).Find(target).Error
}

func (g *jinzhuGorm) Update(ctx context.Context, table string, data map[string]interface{}, whereCondition map[string]interface{}, txs ...interface{}) error {

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

	g.dbMaster = GetGormTracerWithContext(ctx, g.tracer, g.dbMaster)

	if len(txs) > 0 && txs[0] != nil {
		tx, _ := txs[0].(*gorm.DB)
		return tx.Exec(query, vals...).Error
	}
	return g.dbMaster.Exec(query, vals...).Error
}

func (g *jinzhuGorm) Delete(ctx context.Context, table string, whereCondition map[string]interface{}, txs ...interface{}) error {

	query := fmt.Sprintf("DELETE %s WHERE", table)
	whereVals := make([]interface{}, 0)
	whereField := make([]string, 0)

	for k, v := range whereCondition {
		whereVals = append(whereVals, v)
		whereField = append(whereField, k+"=?")
	}

	query = query + strings.Join(whereField, " AND ")
	g.dbMaster = GetGormTracerWithContext(ctx, g.tracer, g.dbMaster)

	if len(txs) > 0 && txs[0] != nil {
		tx, _ := txs[0].(*gorm.DB)
		return tx.Exec(query, whereVals...).Error
	}
	return g.dbMaster.Exec(query, whereVals...).Error
}

func (g *jinzhuGorm) QueryExec(ctx context.Context, txp interface{}, query string, args ...interface{}) error {

	g.dbMaster = GetGormTracerWithContext(ctx, g.tracer, g.dbMaster)

	if txp != nil {
		tx, _ := txp.(*gorm.DB)
		return tx.Exec(query, args...).Error
	}
	return g.dbMaster.Exec(query, args...).Error
}

func (g *jinzhuGorm) QueryRaw(ctx context.Context, txp interface{}, target interface{}, sql string, values ...interface{}) error {

	g.dbMaster = GetGormTracerWithContext(ctx, g.tracer, g.dbMaster)

	if txp != nil {
		tx, _ := txp.(*gorm.DB)
		return tx.Raw(sql, values).Find(target).Error
	}
	return g.dbMaster.Raw(sql, values).Find(target).Error
}
