package server

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type User struct {
	ID       int
	Username string
	Email    string
	Password string
}

func TestAuth() {
	db, err := ConnectDB()
	defer db.Close()

	// Test the signup and login functions.
	user := &User{Username: "john", Email: "john@example.com", Password: "secret"}
	err = Signup(user)
	if err != nil {
		fmt.Println(err)
	}
	loggedInUser, err := Login("john", "secret")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Logged in user: %v\n", loggedInUser)
}

func Signup(user *User) error {
	db, err := ConnectDB()
	defer db.Close()
	// Insert the new user into the database.
	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, user.Password)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok && pgErr.Code == "23505" {
			// Unique constraint violation, the user already exists.
			return fmt.Errorf("user already exists")
		}
		return err
	}
	return nil
}

func Login(username string, password string) (*User, error) {
	db, err := ConnectDB()
	defer db.Close()
	// Find the user with the given username and password.
	user := &User{}
	err = db.QueryRow("SELECT id, username, email, password FROM users WHERE username = $1 AND password = $2", username, password).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid username or password")
		}
		return nil, err
	}
	return user, nil
}
