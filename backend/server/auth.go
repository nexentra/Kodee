package server

import (
	// "crypto/rand"
	"database/sql"
	"net/mail"
	"os"

	// "encoding/base64"
	"fmt"
	// "net/smtp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            int
	Username      string
	Email         string
	Password      string
}

// JWT secret key
var jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// JWT token claims
type jwtClaims struct {
	Username string `json:"username"`
	ID 	 int    `json:"id"`
	jwt.StandardClaims
}

func TestAuth() {
	db, err := ConnectDB()
	defer db.Close()

	// Test the signup and login functions.
	user := &User{Username: "test1", Email: "kaka400068@gmail.com", Password: "test1"}
	token, err := SignUp(user)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Signed up user: %v\n", user)
	fmt.Printf("JWT token: %v\n", token)
	loggedInUser, err := Login("test1", "test1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Logged in user: %v\n", loggedInUser)

	// Test the token validation function.
	userObj, err := Me(loggedInUser)
	fmt.Println("the claims",userObj,err)
	if err != nil {
		fmt.Println("I am the storm",userObj,err)
	}
}

func SignUp(user *User) (string, error) {
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return "", fmt.Errorf("missing required fields")
	}
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return "", fmt.Errorf("invalid email address")
	}
	db, err := ConnectDB()
	defer db.Close()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	// Generate verification token.
	// token := generateVerificationToken(user.Email)
	// Insert the new user into the database.
	// _, err = db.Exec("INSERT INTO users (username, email, password, email_verified, verification_token) VALUES ($1, $2, $3, $4, $5)", user.Username, user.Email, hashedPassword, false, token)
	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, hashedPassword)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok && pgErr.Code == "23505" {
			// Unique constraint violation, the user already exists.
			return "", fmt.Errorf("user already exists")
		}
		return "", err
	}
	// Send verification email to the user.
	// err = sendVerificationEmail(user.Email, token)
	// if err != nil {
	//     return fmt.Errorf("failed to send verification email: %s", err.Error())
	// }
	// Generate a JWT token for the user.
	token, err := generateToken(user.Username, user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %s", err.Error())
	}
	return token, nil
}

func Login(username string, password string) (string, error) {
	if username == "" || password == "" {
		return "", fmt.Errorf("missing required fields")
	}
	db, err := ConnectDB()
	defer db.Close()

	// Find the user with the given username.
	user := &User{}
	err = db.QueryRow("SELECT id, username, email, password FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("invalid username or password")
		}
		return "", err
	}

	// Compare the hashed password with the plain-text password.
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	// Generate a JWT token for the user.
	token, err := generateToken(user.Username, user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %s", err.Error())
	}

	return token, nil
}

func generateToken(username string, id int) (string, error) {
	// Create a new token with claims.
	claims := jwtClaims{
		username,
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			Issuer:    "kodee",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key.
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Me(tokenString string) (User, error) {
	// Parse the JWT token.
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})
	if err != nil {
		return User{}, fmt.Errorf("failed to parse token: %s", err.Error())
	}

	// Verify the token signature and expiration.
	if !token.Valid {
		return User{}, fmt.Errorf("invalid token")
	}

	// Extract the claims from the token.
	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		return User{}, fmt.Errorf("invalid token claims")
	}

	// Retrieve the user ID from the database.
	db, err := ConnectDB()
	if err != nil {
		return User{}, err
	}
	defer db.Close()

	var user User
	// return whole user object if id matches in db
	err = db.QueryRow("SELECT id, username, email FROM users WHERE id = $1 AND username = $2", claims.ID, claims.Username).Scan(&user.ID, &user.Username, &user.Email)

	if err != nil {
		return User{}, err
	}

	return user, nil
}


//for email verification
// func generateVerificationToken(email string) string {
//     randBytes := make([]byte, 32)
//     _, err := rand.Read(randBytes)
//     if err != nil {
//         panic(err)
//     }
//     token := base64.URLEncoding.EncodeToString(randBytes) + email
//     return token
// }

// func sendVerificationEmail(email string, token string) error {
//     // Construct verification link.
//     link := "https://example.com/verify-email?token=" + token
//     // Construct email message.
//     message := "Please click the following link to verify your email address:\n\n" + link
//     // Send email.
//     err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", "example@gmail.com", "yourapppassword", "smtp.gmail.com"), "noreply@example.com", []string{email}, []byte(message))
//     if err != nil {
//         return err
//     }
//     return nil
// }
