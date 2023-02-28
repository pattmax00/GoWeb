package models

import (
	"GoWeb/app"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int64
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GetCurrentUser finds the currently logged-in user by session cookie
func GetCurrentUser(app *app.App, r *http.Request) (User, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Println("Error getting session cookie")
		log.Println(err)
		return User{}, err
	}

	var userId int64

	// Query row by session cookie
	err = app.Db.QueryRow("SELECT \"Id\" FROM public.\"Session\" WHERE \"AuthToken\" = $1", cookie.Value).Scan(&userId)
	if err != nil {
		log.Println("Error querying session row with session: " + cookie.Value)
		return User{}, err
	}

	return GetUserById(app, userId)
}

// GetUserById finds a users table row in the database by id and returns a struct representing this row
func GetUserById(app *app.App, id int64) (User, error) {
	user := User{}

	// Query row by id
	row, err := app.Db.Query("SELECT \"Id\", \"Username\", \"Password\", \"CreatedAt\", \"UpdatedAt\" FROM public.\"User\" WHERE \"Id\" = $1", id)
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

	return user, nil
}

// GetUserByUsername finds a users table row in the database by username and returns a struct representing this row
func GetUserByUsername(app *app.App, username string) (User, error) {
	user := User{}

	// Query row by username
	row, err := app.Db.Query("SELECT \"Id\", \"Username\", \"Password\", \"CreatedAt\", \"UpdatedAt\" FROM public.\"User\" WHERE \"Username\" = $1", username)
	if err != nil {
		log.Println("Error querying user row with username: " + username)
		return User{}, err
	}

	defer func(row *sql.Rows) {
		err := row.Close()
		if err != nil {
			log.Println("Error closing database row")
			log.Println(err)
		}
	}(row)

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

	sqlStatement := "INSERT INTO public.\"User\" (\"Username\", \"Password\", \"CreatedAt\", \"UpdatedAt\") VALUES ($1, $2, $3, $4) RETURNING \"Id\""
	err = app.Db.QueryRow(sqlStatement, username, string(hash), createdAt, updatedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("Error creating user row")
		return User{}, err
	}

	return GetUserById(app, lastInsertId)
}

// AuthenticateUser validates the password for the specified user
func AuthenticateUser(app *app.App, w http.ResponseWriter, username string, password string) (Session, error) {
	var user User

	// Query row by username
	err := app.Db.QueryRow("SELECT \"Id\", \"Username\", \"Password\", \"CreatedAt\", \"UpdatedAt\" FROM public.\"User\" WHERE \"Username\" = $1", username).Scan(&user.Id, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println("Authentication error (user not found) for user:" + username)
		log.Println(err)
		return Session{}, err
	}
	fmt.Println(user)

	// Validate password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil { // Failed to validate password, doesn't match
		log.Println("Authentication error (incorrect password) for user:" + username)
		log.Println(err)
		return Session{}, err
	} else {
		return CreateSession(app, w, user.Id)
	}
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
	err = DeleteSessionByAuthToken(app, w, cookie.Value)
	if err != nil {
		log.Println("Error deleting session by auth token")
		log.Println(err)
		return
	}

	// Delete cookie
	cookie = &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
}
