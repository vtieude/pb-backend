package db_manager

import (
	"context"
	"database/sql"
	"pb-backend/graph/model"
	"time"

	"github.com/elgris/sqrl"
	"github.com/jmoiron/sqlx"
)

type IDb interface {
	Select(ctx context.Context, dest interface{}, sqlizer sqrl.Sqlizer) error
	Get(ctx context.Context, dest interface{}, sqlizer sqrl.Sqlizer) error
	Exec(ctx context.Context, sqlizer sqrl.Sqlizer) (sql.Result, error)
	AddPagination(sq *sqrl.SelectBuilder, pagination *model.Pagination) (*sqrl.SelectBuilder, error)
}
type iSqlxDb interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type DB struct {
	DB iSqlxDb
}

var sqlxDB *DB

func OpenConnectTion() *DB {
	db, err := sqlx.Open("mysql", "root:qweqwe@tcp(127.0.0.1:3307)/app_db?parseTime=true")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	sqlxDB = &DB{DB: db}
	return sqlxDB
}

func (db *DB) Get(ctx context.Context, dest interface{}, sqlizer sqrl.Sqlizer) error {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return err
	}
	sqlxDb := db.DB

	return sqlxDb.GetContext(ctx, dest, query, args...)
}

func (db *DB) Select(ctx context.Context, dest interface{}, sqlizer sqrl.Sqlizer) error {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return err
	}

	sqlxDb := db.DB
	return sqlxDb.SelectContext(ctx, dest, query, args...)
}

func (db *DB) Exec(ctx context.Context, sqlizer sqrl.Sqlizer) (sql.Result, error) {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return nil, err
	}

	sqlxDb := db.DB

	res, err := sqlxDb.ExecContext(ctx, query, args...)
	return res, err
}

func (db *DB) AddPagination(sq *sqrl.SelectBuilder, pagination *model.Pagination) (*sqrl.SelectBuilder, error) {
	if pagination != nil {
		if pagination.PerPage != nil && pagination.Page != nil {
			offset := uint64((*pagination.Page - 1) * *pagination.PerPage)
			sq = sq.Offset(offset).Limit(uint64(*pagination.PerPage))
		}
		if pagination.Sort != nil {
			sq = sq.OrderBy(pagination.Sort...)
		}
	}
	return sq, nil
}
