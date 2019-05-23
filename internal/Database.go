package internal

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "ec2-54-225-242-183.compute-1.amazonaws.com"
	port     = 5432
	user     = "aimzpnysofwypw"
	password = "de56c756197c4d8f41745acf76ff3df6c3cc39852c7eb5572d173778d7ba28de"
	dbname   = "dbif64ksnitjje"
)

func ConexaoPostgres() (db *sql.DB, err error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db, err
}
