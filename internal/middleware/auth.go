package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	cookieName = "user_id"
	secretKey  = "92KbwrTL3zWMqD7egj4L5Y7"
)

func MakeAuthMiddleware(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// FIXME: do better exceptions
			if r.URL.Path == "/api/user/urls" {
				next.ServeHTTP(w, r)
				return
			}

			cookie, err := r.Cookie(cookieName)
			if err != nil || !validateCookie(cookie.Value) {
				userID := uuid.New().String()
				signedValue := SignCookie(userID)
				http.SetCookie(w, &http.Cookie{
					Name:    cookieName,
					Value:   signedValue,
					Expires: time.Now().Add(24 * time.Hour),
					Path:    "/",
				})
				r.AddCookie(&http.Cookie{Name: cookieName, Value: signedValue})
			}
			next.ServeHTTP(w, r)
		})
	}
}

func SignCookie(value string) string {
	// FIXME: move to jwt
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(value))
	signature := hex.EncodeToString(h.Sum(nil))
	return value + "." + signature
}

func validateCookie(signedValue string) bool {
	parts := strings.Split(signedValue, ".")
	if len(parts) != 2 {
		return false
	}

	return SignCookie(parts[0]) == signedValue
}

func GetUserID(r *http.Request) (string, bool) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		return "", false
	}
	parts := strings.Split(cookie.Value, ".")
	if len(parts) != 2 {
		return "", false
	}
	return parts[0], true
}
