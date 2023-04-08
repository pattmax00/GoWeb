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
	Id         int64
	UserId     int64
	AuthToken  string
	RememberMe bool
	CreatedAt  time.Time
}

const sessionColumnsNoId = "\"UserId\", \"AuthToken\",\"RememberMe\", \"CreatedAt\""
const sessionColumns = "\"Id\", " + sessionColumnsNoId
const sessionTable = "public.\"Session\""

const (
	selectSessionByAuthToken      = "SELECT " + sessionColumns + " FROM " + sessionTable + " WHERE \"AuthToken\" = $1"
	selectAuthTokenIfExists       = "SELECT EXISTS(SELECT 1 FROM " + sessionTable + " WHERE \"AuthToken\" = $1)"
	insertSession                 = "INSERT INTO " + sessionTable + " (" + sessionColumnsNoId + ") VALUES ($1, $2, $3, $4) RETURNING \"Id\""
	deleteSessionByAuthToken      = "DELETE FROM " + sessionTable + " WHERE \"AuthToken\" = $1"
	deleteSessionsOlderThan30Days = "DELETE FROM " + sessionTable + " WHERE \"CreatedAt\" < NOW() - INTERVAL '30 days'"
	deleteSessionsOlderThan6Hours = "DELETE FROM " + sessionTable + " WHERE \"CreatedAt\" < NOW() - INTERVAL '6 hours' AND \"RememberMe\" = false"
)

// CreateSession creates a new session for a user
func CreateSession(app *app.App, w http.ResponseWriter, userId int64, remember bool) (Session, error) {
	session := Session{}
	session.UserId = userId
	session.AuthToken = generateAuthToken(app)
	session.RememberMe = remember
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
		return CreateSession(app, w, userId, remember)
	}

	// Insert session into database
	err = app.Db.QueryRow(insertSession, session.UserId, session.AuthToken, session.RememberMe, session.CreatedAt).Scan(&session.Id)
	if err != nil {
		log.Println("Error inserting session into database")
		return Session{}, err
	}

	createSessionCookie(app, w, session)
	return session, nil
}

func GetSessionByAuthToken(app *app.App, authToken string) (Session, error) {
	session := Session{}

	err := app.Db.QueryRow(selectSessionByAuthToken, authToken).Scan(&session.Id, &session.UserId, &session.AuthToken, &session.RememberMe, &session.CreatedAt)
	if err != nil {
		log.Println("Error getting session by auth token")
		return Session{}, err
	}

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
	cookie := &http.Cookie{}
	if session.RememberMe {
		cookie = &http.Cookie{
			Name:     "session",
			Value:    session.AuthToken,
			Path:     "/",
			MaxAge:   2592000 * 1000, // 30 days in ms
			HttpOnly: true,
			Secure:   true,
		}
	} else {
		cookie = &http.Cookie{
			Name:     "session",
			Value:    session.AuthToken,
			Path:     "/",
			MaxAge:   21600 * 1000, // 6 hours in ms
			HttpOnly: true,
			Secure:   true,
		}
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

// ScheduledSessionCleanup deletes expired sessions from the database
func ScheduledSessionCleanup(app *app.App) {
	// Delete sessions older than 30 days (remember me sessions)
	_, err := app.Db.Exec(deleteSessionsOlderThan30Days)
	if err != nil {
		log.Println("Error deleting 30 day expired sessions from database")
		log.Println(err)
	}

	// Delete sessions older than 6 hours
	_, err = app.Db.Exec(deleteSessionsOlderThan6Hours)
	if err != nil {
		log.Println("Error deleting 6 hour expired sessions from database")
		log.Println(err)
	}

	log.Println("Deleted expired sessions from database")
}
