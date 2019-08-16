package driver

import (
	"database/sql"

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

func ConnectSQL(url string) (*DB, error) {

	db, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}
	dbConn.SQL = db
	return dbConn, err
}
