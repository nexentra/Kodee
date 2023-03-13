package server

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"
)

type PortfolioTables struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	UDEF        json.RawMessage `json:"udef"`
	User_id     int             `json:"user_id"`
	User        User            `json:"user"`
}

func TestOut() {
	// Open a connection to the database

	db, err := ConnectDB()
	// Create a new person
	p1 := PortfolioTables{Name: "Alice2", Description: "Hello", UDEF: json.RawMessage(`{"test": "test"}`), User_id: 1}
	err = CreatePortfolioTables(&p1)
	if err != nil {
		fmt.Println(err)
	}

	// Retrieve a person by ID
	p2, err := GetPortfolioTables(p1.ID, p1.User_id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p2)

	// Update a person's age
	p1.Name = "Bob"
	err = UpdatePortfolioTables(&p1)
	if err != nil {
		fmt.Println(err)
	}

	// Delete a person by ID
	// err = DeletePortfolioTables( p1.ID)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// Verify that the person has been deleted
	p3, err := GetPortfolioTables(p1.ID, p1.User_id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p3)

	defer db.Close()
}

func CreatePortfolioTables(p *PortfolioTables) error {
	db, err := ConnectDB()
	sqlStatement := `
	INSERT INTO portfolio (name, description, udef, user_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id`
	id := 0
	err = db.QueryRow(sqlStatement, p.Name, p.Description, p.UDEF, p.User_id).Scan(&id)
	if err != nil {
		return err
	}
	p.ID = id
	defer db.Close()
	return nil
}

func GetPortfolioTables(id int, user_id int) (PortfolioTables, error) {
	db, err := ConnectDB()
	sqlStatement := `SELECT * FROM portfolio WHERE id=$1 AND user_id=$2`
	row := db.QueryRow(sqlStatement, id, user_id)
	var p PortfolioTables
	err = row.Scan(&p.ID, &p.Name, &p.Description, &p.UDEF, &p.User_id)
	defer db.Close()
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return p, nil
	case nil:
		return p, nil
	default:
		fmt.Println(err)
		return p, err
	}
}

func UpdatePortfolioTables(p *PortfolioTables) error {
	db, err := ConnectDB()
	sqlStatement := `
	UPDATE portfolio
	SET name = $2, description = $3, udef = $4
	WHERE id = $1 AND user_id = $5`
	_, err = db.Exec(sqlStatement, p.ID, p.Name, p.Description, p.UDEF, p.User_id)
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

func DeletePortfolioTables(id int, user_id int) error {
	db, err := ConnectDB()
	sqlStatement := `DELETE FROM portfolio WHERE id=$1 AND user_id=$2`
	_, err = db.Exec(sqlStatement, id, user_id)
	if err != nil {
		return err
	}

	defer db.Close()
	return nil
}
