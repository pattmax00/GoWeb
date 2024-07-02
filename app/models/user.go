package models

import (
	"GoWeb/app"
	"crypto/sha256"
	"encoding/hex"
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

// CurrentUser finds the currently logged-in user by session cookie
func CurrentUser(app *app.Deps, r *http.Request) (User, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return User{}, err
	}

	session, err := SessionByAuthToken(app, cookie.Value)
	if err != nil {
		return User{}, err
	}

	return UserById(app, session.UserId)
}

// UserById finds a User table row in the database by id and returns a struct representing this row
func UserById(app *app.Deps, id int64) (User, error) {
	user := User{}

	err := app.Db.QueryRow(selectUserById, id).Scan(&user.Id, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// UserByUsername finds a User table row in the database by username and returns a struct representing this row
func UserByUsername(app *app.Deps, username string) (User, error) {
	user := User{}

	err := app.Db.QueryRow(selectUserByUsername, username).Scan(&user.Id, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// CreateUser creates a User table row in the database
func CreateUser(app *app.Deps, username string, password string, createdAt time.Time, updatedAt time.Time) (User, error) {
	// Get sha256 hash of password then get bcrypt hash to store
	hash256 := sha256.New()
	hash256.Write([]byte(password))
	hashSum := hash256.Sum(nil)
	hashString := hex.EncodeToString(hashSum)
	hash, err := bcrypt.GenerateFromPassword([]byte(hashString), bcrypt.DefaultCost)
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

	return UserById(app, lastInsertId)
}

// AuthenticateUser validates the password for the specified user
func AuthenticateUser(app *app.Deps, w http.ResponseWriter, username string, password string, remember bool) (Session, error) {
	var user User

	err := app.Db.QueryRow(selectUserByUsername, username).Scan(&user.Id, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		slog.Info("user not found: " + username)
		return Session{}, err
	}

	// Get sha256 hash of password then check bcrypt
	hash256 := sha256.New()
	hash256.Write([]byte(password))
	hashSum := hash256.Sum(nil)
	hashString := hex.EncodeToString(hashSum)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(hashString))
	if err != nil { // Failed to validate password, doesn't match
		slog.Info("incorrect password:" + username)
		return Session{}, err
	} else {
		return CreateSession(app, w, user.Id, remember)
	}
}

// LogoutUser deletes the session cookie and AuthToken from the database
func LogoutUser(app *app.Deps, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return
	}

	err = DeleteSessionByAuthToken(app, w, cookie.Value)
	if err != nil {
		return
	}
}
