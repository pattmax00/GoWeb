package models

import (
	"GoWeb/app"
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/gob"
	"encoding/hex"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int64
	Username  string
	Password  string
	CreatedAt string
	UpdatedAt string
}

// GetUserById finds a users table row in the database by id and returns a struct representing this row
func GetUserById(app *app.App, id int64) (User, error) {
	user := User{}

	// Query row by id
	row, err := app.Db.Query("SELECT id, username, password, created_at, updated_at FROM users WHERE id = $1", id)
	if err != nil {
		log.Println("Error querying user row with id: " + strconv.FormatInt(id, 10))
		return User{}, err
	}

	defer func(row *sql.Rows) {
		err := row.Close()
		if err != nil {
			log.Println("Error closing database row")
			log.Println(err)
		}
	}(row)

	// Feed row data into user struct
	row.Next()
	err = row.Scan(&user.Id, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println("Error reading queried row from database")
		log.Println(err)
		return User{}, err
	}

	return user, nil
}

// CreateUser creates a users table row in the database
func CreateUser(app *app.App, username string, password string, createdAt time.Time, updatedAt time.Time) (User, error) {
	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password when creating user")
		return User{}, err
	}

	var lastInsertId int64

	sqlStatement := "INSERT INTO users (username, password, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id"
	err = app.Db.QueryRow(sqlStatement, username, string(hash), createdAt, updatedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("Error creating user row")
		log.Println(err)
		return User{}, err
	}

	return GetUserById(app, lastInsertId)
}

// AuthenticateUser validates the password for the specified user if it matches a session cookie is created and returned
func AuthenticateUser(app *app.App, w http.ResponseWriter, username string, password string) (string, error) {
	var hashedPassword []byte

	// Query row by username, scan password column
	err := app.Db.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&hashedPassword)
	if err != nil {
		log.Println("Unable to find row with username: " + username)
		log.Println(err)
		return "", err
	}

	// Convert passed password to []byte
	passwordBuffer := &bytes.Buffer{}
	err = gob.NewEncoder(passwordBuffer).Encode(password)
	if err != nil {
		log.Println("Unable to convert passed password to []byte")
		log.Println(err)
		return "", err
	}
	passwordByteSlice := passwordBuffer.Bytes()

	// Validate password
	err = bcrypt.CompareHashAndPassword(hashedPassword, passwordByteSlice)
	if err != nil { // Failed to validate password, doesn't match
		log.Println("Authentication error (incorrect password) for user:" + username)
		log.Println(err)
		return "", err
	} else {
		return createSessionCookie(app, w, username)
	}
}

// createSessionCookie creates a new session token and cookie and returns the token value
func createSessionCookie(app *app.App, w http.ResponseWriter, username string) (string, error) {
	// Generate random 64 character string (alpha-numeric)
	buff := make([]byte, int(math.Ceil(float64(64)/2)))
	_, err := rand.Read(buff)
	if err != nil {
		log.Println("Error creating random buffer for token value")
		log.Println(err)
		return "", err
	}
	str := hex.EncodeToString(buff)
	token := str[:64]

	// Store token in auth_token column of users table
	sqlStatement := "UPDATE users SET auth_token = $1 WHERE username = $2"
	err = app.Db.QueryRow(sqlStatement, token, username).Scan()
	if err != nil {
		log.Println("Error setting auth_token column in users table")
		log.Println(err)
		return "", err
	}

	// Create session cookie, containing token
	cookie := &http.Cookie{
		Name:   "session",
		Value:  token,
		Path:   "/",
		MaxAge: 86400,
	}

	http.SetCookie(w, cookie)

	return token, nil
}

// ValidateSessionCookie validates the session cookie and returns the username of the user if valid
func ValidateSessionCookie(app *app.App, r *http.Request) (string, error) {
	// Get cookie from request
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Println("Error getting cookie from request")
		log.Println(err)
		return "", err
	}

	// Query row by token
	var username string
	err = app.Db.QueryRow("SELECT username FROM users WHERE auth_token = $1", cookie.Value).Scan(&username)
	if err != nil {
		log.Println("Error querying row by token")
		log.Println(err)
		return "", err
	}

	return username, nil
}

// LogoutUser deletes the session cookie and token from the database
func LogoutUser(app *app.App, w http.ResponseWriter, r *http.Request) {
	// Get cookie from request
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Println("Error getting cookie from request")
		log.Println(err)
		return
	}

	// Set token to empty string
	sqlStatement := "UPDATE users SET auth_token = $1 WHERE auth_token = $2"
	err = app.Db.QueryRow(sqlStatement, "", cookie.Value).Scan()
	if err != nil {
		log.Println("Error setting auth_token column in users table")
		log.Println(err)
		return
	}

	// Delete cookie
	cookie = &http.Cookie{
		Name:   "",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
}
