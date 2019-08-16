package driver

import (
	"database/sql"
	"fmt"

	//"github.com/CezarGarrido/sqllogs"
	_ "github.com/lib/pq"
)

// DB ...
type DB struct {
	SQL *sql.DB
	// Mgo *mgo.database
}

// DBConn ...
var dbConn = &DB{}

func ConnectSQL(host, port, user, password, dbname string) (*DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	dbConn.SQL = db
	return dbConn, err
}
