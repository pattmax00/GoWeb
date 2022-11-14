package security

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"math"
	"net/http"
	"time"
)

// GenerateCsrfToken generates a csrf token and assigns it to a cookie for double submit cookie csrf protection
func GenerateCsrfToken(w http.ResponseWriter, r *http.Request) (string, error) {
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

	// Create session cookie, containing token
	cookie := &http.Cookie{
		Name:     "csrf",
		Value:    token,
		Path:     "/",
		MaxAge:   1800,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, cookie)

	return token, nil
}

// VerifyCsrfToken verifies the csrf token
func VerifyCsrfToken(r *http.Request) (bool, error) {
	// Get csrf cookie
	cookie, err := r.Cookie("csrf")
	if err != nil {
		log.Println("Error getting csrf cookie")
		log.Println(err)
		return false, err
	}

	// Get csrf token from form
	token := r.FormValue("csrf")

	// Compare csrf cookie and csrf token
	if cookie.Value == token {
		return true, nil
	}

	return false, nil
}
