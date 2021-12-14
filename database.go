package main

import (
	"database/sql"
	"fmt"
)

func OpenDbConnection(cfg Config) (*sql.DB, error) {
	// Open up our database connection.
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.DbUser, cfg.DbPsw, cfg.DbHost, cfg.DbName)
	db, err := sql.Open("mysql", connectionString)

	// if there is an error opening the connection, handle it
	if err != nil {
		println("cannot connect db", err, connectionString)
		return db, err
	}

	// defer the close till after the main function has finished
	// executing
	return db, nil
}
