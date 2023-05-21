package server

import (
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
	p1 := PortfolioTables{Name: "jsontest", Description: "test", UDEF: json.RawMessage(`
	{
		"tableName": "firstTable",
		"tableDescription": "firstTableDescription is this",
		"tableColumnNames": ["columnName", "columnDescription", "columnType"],
		"tableColumns": [
			{
				"columnName": "firstColumn",
				"columnDescription": "firstColumnDescription is this",
				"columnType": "string"
			},
			{
				"columnName": "secondColumn",
				"columnDescription": "secondColumnDescription is this",
				"columnType": "string"

			},
			{
				"columnName": "thirdColumn",
				"columnDescription": "thirdColumnDescription is this",
				"columnType": "string"
			}
		]
	}
	`), User_id: 1}
	err = CreatePortfolioTables("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3QxIiwiaWQiOjEsImV4cCI6MTY3OTQzMDU0MCwiaXNzIjoia29kZWUifQ.H58gRmNpDh_wakDxfOAlJFM9YotjvN2VjMuRkmWgmLc",&p1)
	if err != nil {
		fmt.Println(err)
	}

	// Retrieve a person by ID
	p2, err := GetPortfolioTables("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3QxIiwiaWQiOjEsImV4cCI6MTY3OTQzMDU0MCwiaXNzIjoia29kZWUifQ.H58gRmNpDh_wakDxfOAlJFM9YotjvN2VjMuRkmWgmLc",p1.ID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p2)

	// Update a person's age
	p1.Name = "Bob"
	err = UpdatePortfolioTables("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3QxIiwiaWQiOjEsImV4cCI6MTY3OTQzMDU0MCwiaXNzIjoia29kZWUifQ.H58gRmNpDh_wakDxfOAlJFM9YotjvN2VjMuRkmWgmLc",&p1)
	if err != nil {
		fmt.Println(err)
	}

	// Delete a person by ID
	// err = DeletePortfolioTables( p1.ID)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// Verify that the person has been deleted
	p3, err := GetPortfolioTables("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3QxIiwiaWQiOjEsImV4cCI6MTY3OTQzMDU0MCwiaXNzIjoia29kZWUifQ.H58gRmNpDh_wakDxfOAlJFM9YotjvN2VjMuRkmWgmLc",p1.ID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p3)

	p4, err := GetSinglePortfolio("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3QxIiwiaWQiOjEsImV4cCI6MTY3OTQzMDU0MCwiaXNzIjoia29kZWUifQ.H58gRmNpDh_wakDxfOAlJFM9YotjvN2VjMuRkmWgmLc",p1.ID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p4)

	defer db.Close()
}

func CreatePortfolioTables(token string, p *PortfolioTables) error {
	user, err := Me(token)
	if err != nil {
		return fmt.Errorf("Please login to continue!!")
	}
	db, err := ConnectDB()
	defer db.Close()
	sqlStatement := `
	INSERT INTO portfolio (name, description, udef, user_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id`
	id := 0
	err = db.QueryRow(sqlStatement, p.Name, p.Description, p.UDEF, user.ID).Scan(&id)
	if err != nil {
		return err
	}
	p.ID = id
	return nil
}

func GetPortfolioTables(token string, id int) ([]PortfolioTables, error) {
	user, err := Me(token)
	if err != nil {
		return []PortfolioTables{}, fmt.Errorf("Please login to continue!!")
	}
	db, err := ConnectDB()
	defer db.Close()
	sqlStatement := `SELECT * FROM portfolio WHERE user_id=$1`
	rows, err := db.Query(sqlStatement, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var portfolios []PortfolioTables
	for rows.Next() {
		var p PortfolioTables
		err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.UDEF, &p.User_id)
		if err != nil {
			return nil, err
		}
		portfolios = append(portfolios, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return portfolios, nil
}


func GetSinglePortfolio(token string, portfolioID int) (PortfolioTables, error) {
	user, err := Me(token)
	if err != nil {
		return PortfolioTables{}, fmt.Errorf("Please login to continue!!")
	}

	db, err := ConnectDB()
	if err != nil {
		return PortfolioTables{}, err
	}
	defer db.Close()

	var portfolio PortfolioTables
	err = db.QueryRow("SELECT * FROM portfolio WHERE id=$1 AND user_id=$2", portfolioID, user.ID).Scan(
		&portfolio.ID,
		&portfolio.Name,
		&portfolio.Description,
		&portfolio.UDEF,
		&portfolio.User_id,
	)

	if err != nil {
		return PortfolioTables{}, err
	}

	return portfolio, nil
}


func UpdatePortfolioTables(token string, p *PortfolioTables) error {
	_, err := Me(token)
	if err != nil {
		return fmt.Errorf("Please login to continue!!")
	}
	db, err := ConnectDB()
	defer db.Close()
	sqlStatement := `
	UPDATE portfolio
	SET name = $2, description = $3, udef = $4
	WHERE id = $1 AND user_id = $5`
	_, err = db.Exec(sqlStatement, p.ID, p.Name, p.Description, p.UDEF, p.User_id)
	if err != nil {
		return err
	}
	return nil
}

func DeletePortfolioTables(token string, id int, user_id int) error {
	_, err := Me(token)
	if err != nil {
		return fmt.Errorf("Please login to continue!!")
	}
	db, err := ConnectDB()
	defer db.Close()
	sqlStatement := `DELETE FROM portfolio WHERE id=$1 AND user_id=$2`
	_, err = db.Exec(sqlStatement, id, user_id)
	if err != nil {
		return err
	}

	return nil
}
