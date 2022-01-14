package entities

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"pb-backend/graph/model"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/elgris/sqrl"
	"github.com/jmoiron/sqlx"
)

const ConfigKey = "pbconfig"

type PbConfig struct {
	DbUser     string `yaml:"DbUser"`
	DbPsw      string `yaml:"DbPsw"`
	DbName     string `yaml:"DbName"`
	DbHost     string `yaml:"DbHost"`
	AppPort    string `yaml:"AppPort"`
	AppProdEnv bool   `yaml:"AppProdEnv"`
}

type MyCustomClaims struct {
	Username string `json:"username"`
	ID       int    `json:"userid"`
	jwt.StandardClaims
}

// type IDb interface {
// 	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
// 	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
// 	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
// 	AddPagination(sq *sqrl.SelectBuilder, pagination *model.Pagination) (*sqrl.SelectBuilder, error)
// }

type iSqlxDb interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type DBConnection struct {
	DB  iSqlxDb
	log log.Logger
}

var sqlxDB *DBConnection

func OpenConnection(ctx context.Context, log log.Logger) *DBConnection {
	var cfg PbConfig
	cfg, ok := ctx.Value(ConfigKey).(PbConfig)
	if ok {
		dbConString := ""
		if cfg.AppProdEnv {
			dbConString = os.Getenv("DATABASE_URL")
			if dbConString == "" {
				log.Fatal("$DATABASE_URL is not set")
			}

		} else {
			dbConString = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", cfg.DbUser, cfg.DbPsw, cfg.DbHost, cfg.DbName)
		}
		db, err := sqlx.Open("mysql", dbConString)
		if err != nil {
			panic(err)
		}
		// See "Important settings" section.
		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(0)
		sqlxDB = &DBConnection{DB: db, log: log}
		return sqlxDB
	}
	db, err := sqlx.Open("mysql", "root:qweqwe@tcp(127.0.0.1:3307)/app_db?parseTime=true")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	sqlxDB = &DBConnection{DB: db}
	return sqlxDB
}

func (db *DBConnection) QueryRowContext(ctx context.Context, dest interface{}, sqlizer sqrl.Sqlizer, args ...interface{}) error {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return err
	}
	sqlxDb := db.DB
	db.log.Println(query, args)
	return sqlxDb.GetContext(ctx, dest, query, args...)
}

func (db *DBConnection) QueryContext(ctx context.Context, dest interface{}, sqlizer sqrl.Sqlizer, args ...interface{}) error {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return err
	}
	db.log.Println(query, args)
	sqlxDb := db.DB
	return sqlxDb.SelectContext(ctx, dest, query, args...)
}

func (db *DBConnection) ExecSqrlContext(ctx context.Context, sqlizer sqrl.Sqlizer, args ...interface{}) (sql.Result, error) {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return nil, err
	}
	db.log.Println(query, args)
	sqlxDb := db.DB

	res, err := sqlxDb.ExecContext(ctx, query, args...)
	return res, err
}

func (db *DBConnection) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	sqlxDb := db.DB
	res, err := sqlxDb.ExecContext(ctx, query, args...)
	db.log.Println(query, args)
	return res, err
}

func (db *DBConnection) AddPagination(sq *sqrl.SelectBuilder, pagination *model.Pagination) (*sqrl.SelectBuilder, error) {
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
