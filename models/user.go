package models

import (
	"GoWeb/app"
	"log/slog"
	"net/http"
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

const userColumnsNoId = "\"Username\", \"Password\", \"CreatedAt\", \"UpdatedAt\""
const userColumns = "\"Id\", " + userColumnsNoId
const userTable = "public.\"User\""

const (
	selectUserById       = "SELECT " + userColumns + " FROM " + userTable + " WHERE \"Id\" = $1"
	selectUserByUsername = "SELECT " + userColumns + " FROM " + userTable + " WHERE \"Username\" = $1"
	insertUser           = "INSERT INTO " + userTable + " (" + userColumnsNoId + ") VALUES ($1, $2, $3, $4) RETURNING \"Id\""
)

// GetCurrentUser finds the currently logged-in user by session cookie
func GetCurrentUser(app *app.App, r *http.Request) (User, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return User{}, err
	}

	session, err := GetSessionByAuthToken(app, cookie.Value)
	if err != nil {
		return User{}, err
	}

	return GetUserById(app, session.UserId)
}

// GetUserById finds a User table row in the database by id and returns a struct representing this row
func GetUserById(app *app.App, id int64) (User, error) {
	user := User{}

	err := app.Db.QueryRow(selectUserById, id).Scan(&user.Id, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// GetUserByUsername finds a User table row in the database by username and returns a struct representing this row
func GetUserByUsername(app *app.App, username string) (User, error) {
	user := User{}

	err := app.Db.QueryRow(selectUserByUsername, username).Scan(&user.Id, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// CreateUser creates a User table row in the database
func CreateUser(app *app.App, username string, password string, createdAt time.Time, updatedAt time.Time) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("error hashing password: " + err.Error())
		return User{}, err
	}

	var lastInsertId int64

	err = app.Db.QueryRow(insertUser, username, string(hash), createdAt, updatedAt).Scan(&lastInsertId)
	if err != nil {
		slog.Error("error creating user row: " + err.Error())
		return User{}, err
	}

	return GetUserById(app, lastInsertId)
}

// AuthenticateUser validates the password for the specified user
func AuthenticateUser(app *app.App, w http.ResponseWriter, username string, password string, remember bool) (Session, error) {
	var user User

	err := app.Db.QueryRow(selectUserByUsername, username).Scan(&user.Id, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		slog.Info("user not found: " + username)
		return Session{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil { // Failed to validate password, doesn't match
		slog.Info("incorrect password:" + username)
		return Session{}, err
	} else {
		return CreateSession(app, w, user.Id, remember)
	}
}

// LogoutUser deletes the session cookie and AuthToken from the database
func LogoutUser(app *app.App, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return
	}

	err = DeleteSessionByAuthToken(app, w, cookie.Value)
	if err != nil {
		return
	}
}
