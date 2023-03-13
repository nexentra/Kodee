package server

import (
	"database/sql"
	"fmt"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "mydatabase"
)

func ConnectDB() (*sql.DB, error) {
	connectionString := "postgres://ejgtvkwf:ZgfC8vmyY51c4Irs2NNwkfhf32Wq4-am@salt.db.elephantsql.com/ejgtvkwf"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println(err)
	}

	sqlStatement := `
		CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, username TEXT UNIQUE, email TEXT UNIQUE, password TEXT)
		`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		fmt.Println(err)
	}

	// Create a table
	sqlStatement = `
		CREATE TABLE IF NOT EXISTS portfolio (
			id SERIAL PRIMARY KEY,
			name TEXT,
			description TEXT,
			udef JSON,
			user_id INT REFERENCES users(id) ON DELETE CASCADE
		)
		`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		fmt.Println(err)
	}

	return db, err
}
