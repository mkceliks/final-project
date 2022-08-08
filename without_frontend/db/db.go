package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // side effect of the import. Together with the import of the driver, it will be used to open a connection to the database.
)

const ( // constants for the database
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "159987741"
	dbname   = "projectdb"
)

var DB *sql.DB // global variable for the database connection

func init() { // function to initialize the database connection
	var err error
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname) // it is a connection string
	DB, err = sql.Open("postgres", conn)                                                                                     // open a connection to the database
	if err != nil {
		log.Fatal(err)
	}
}
