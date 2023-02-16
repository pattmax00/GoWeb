package security

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"math"
	"net/http"
)

// GenerateCsrfToken generates a csrf token and assigns it to a cookie for double submit cookie csrf protection
func GenerateCsrfToken(w http.ResponseWriter, _ *http.Request) (string, error) {
	// Generate random 64 character string (alpha-numeric)
	buff := make([]byte, int(math.Ceil(float64(64)/2)))
	_, err := rand.Read(buff)
	if err != nil {
		log.Println("Error creating random buffer for csrf token value")
		log.Println(err)
		return "", err
	}
	str := hex.EncodeToString(buff)
	token := str[:64]

	cookie := &http.Cookie{
		Name:     "csrf_token",
		Value:    token,
		Path:     "/",
		MaxAge:   1800,
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, cookie)

	return token, nil
}

// VerifyCsrfToken verifies the csrf token
func VerifyCsrfToken(r *http.Request) (bool, error) {
	cookie, err := r.Cookie("csrf_token")
	if err != nil {
		log.Println("Error getting csrf_token cookie")
		log.Println(err)
		return false, err
	}

	token := r.FormValue("csrf_token")

	if cookie.Value == token {
		return true, nil
	}

	return false, nil
}
