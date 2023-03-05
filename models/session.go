package models

import (
	"GoWeb/app"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"
)

type Session struct {
	Id        int64
	UserId    int64
	AuthToken string
	CreatedAt time.Time
}

const sessionColumnsNoId = "\"UserId\", \"AuthToken\", \"CreatedAt\""
const sessionColumns = "\"Id\", " + sessionColumnsNoId
const sessionTable = "public.\"Session\""

const (
	selectSessionByAuthToken = "SELECT " + sessionColumns + " FROM " + sessionTable + " WHERE \"AuthToken\" = $1"
	selectAuthTokenIfExists  = "SELECT EXISTS(SELECT 1 FROM " + sessionTable + " WHERE \"AuthToken\" = $1)"
	insertSession            = "INSERT INTO " + sessionTable + " (" + sessionColumnsNoId + ") VALUES ($1, $2, $3) RETURNING \"Id\""
	deleteSessionByAuthToken = "DELETE FROM " + sessionTable + " WHERE \"AuthToken\" = $1"
)

// CreateSession creates a new session for a user
func CreateSession(app *app.App, w http.ResponseWriter, userId int64) (Session, error) {
	session := Session{}
	session.UserId = userId
	session.AuthToken = generateAuthToken(app)
	session.CreatedAt = time.Now()

	// If the AuthToken column for any user matches the token, set existingAuthToken to true
	var existingAuthToken bool
	err := app.Db.QueryRow(selectAuthTokenIfExists, session.AuthToken).Scan(&existingAuthToken)
	if err != nil {
		log.Println("Error checking for existing auth token")
		log.Println(err)
		return Session{}, err
	}

	// If duplicate token found, recursively call function until unique token is generated
	if existingAuthToken == true {
		log.Println("Duplicate token found in sessions table, generating new token...")
		return CreateSession(app, w, userId)
	}

	// Insert session into database
	err = app.Db.QueryRow(insertSession, session.UserId, session.AuthToken, session.CreatedAt).Scan(&session.Id)
	if err != nil {
		log.Println("Error inserting session into database")
		return Session{}, err
	}

	createSessionCookie(app, w, session)
	return session, nil
}

// Generates a random 64-byte string
func generateAuthToken(app *app.App) string {
	// Generate random bytes
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		log.Println("Error generating random bytes")
	}

	// Convert random bytes to hex string
	return hex.EncodeToString(b)
}

// createSessionCookie creates a new session cookie
func createSessionCookie(app *app.App, w http.ResponseWriter, session Session) {
	cookie := &http.Cookie{
		Name:     "session",
		Value:    session.AuthToken,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, cookie)
}

// deleteSessionCookie deletes the session cookie
func deleteSessionCookie(app *app.App, w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
}

// DeleteSessionByAuthToken deletes a session from the database by AuthToken
func DeleteSessionByAuthToken(app *app.App, w http.ResponseWriter, authToken string) error {
	// Delete session from database
	_, err := app.Db.Exec(deleteSessionByAuthToken, authToken)
	if err != nil {
		log.Println("Error deleting session from database")
		return err
	}

	deleteSessionCookie(app, w)

	return nil
}
