package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // side effect of the import. Together with the import of the driver, it will be used to open a connection to the database.
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "159987741"
	dbname   = "projectdb"
)

var DB *sql.DB

func init() {
	var err error
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname) // it is a connection string
	DB, err = sql.Open("postgres", conn)                                                                                     // open the connection
	if err != nil {
		log.Fatal(err)
	}
}
