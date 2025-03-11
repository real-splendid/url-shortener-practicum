package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestAuthMiddleware(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()

	t.Run("new user without cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handlerCalled := false
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			cookie, err := r.Cookie(cookieName)
			require.NoError(t, err)
			assert.NotEmpty(t, cookie.Value)
			assert.True(t, validateCookie(cookie.Value))
		})

		handler := MakeAuthMiddleware(logger)(testHandler)
		handler.ServeHTTP(rr, req)

		assert.True(t, handlerCalled)

		resp := rr.Result()
		defer resp.Body.Close()

		cookies := resp.Cookies()
		require.NotEmpty(t, cookies)
		var userIDCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == cookieName {
				userIDCookie = cookie
				break
			}
		}
		require.NotNil(t, userIDCookie)
		assert.NotEmpty(t, userIDCookie.Value)
		assert.True(t, validateCookie(userIDCookie.Value))
	})

	t.Run("existing user with valid cookie", func(t *testing.T) {
		userID := "test-user-id"
		signedValue := SignCookie(userID)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: cookieName, Value: signedValue})

		rr := httptest.NewRecorder()

		handlerCalled := false
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			cookie, err := r.Cookie(cookieName)
			require.NoError(t, err)
			assert.Equal(t, signedValue, cookie.Value)
		})

		handler := MakeAuthMiddleware(logger)(testHandler)
		handler.ServeHTTP(rr, req)

		assert.True(t, handlerCalled)

		resp := rr.Result()
		defer resp.Body.Close()

		cookies := resp.Cookies()
		for _, cookie := range cookies {
			if cookie.Name == cookieName {
				t.Error("No new cookie should be set for existing valid user")
			}
		}
	})

	t.Run("user with invalid cookie", func(t *testing.T) {
		invalidValue := "invalid-user-id.invalid-signature"

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: cookieName, Value: invalidValue})

		rr := httptest.NewRecorder()

		handlerCalled := false
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
		})

		handler := MakeAuthMiddleware(logger)(testHandler)
		handler.ServeHTTP(rr, req)

		assert.True(t, handlerCalled)

		resp := rr.Result()
		defer resp.Body.Close()

		cookies := resp.Cookies()
		require.NotEmpty(t, cookies)
		var userIDCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == cookieName {
				userIDCookie = cookie
				break
			}
		}
		require.NotNil(t, userIDCookie)
		assert.NotEqual(t, invalidValue, userIDCookie.Value)
		assert.True(t, validateCookie(userIDCookie.Value))
	})

	t.Run("api user urls path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
		rr := httptest.NewRecorder()

		handlerCalled := false
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
		})

		handler := MakeAuthMiddleware(logger)(testHandler)
		handler.ServeHTTP(rr, req)

		assert.True(t, handlerCalled)

		resp := rr.Result()
		defer resp.Body.Close()
	})
}

func TestSignCookie(t *testing.T) {
	t.Run("sign cookie", func(t *testing.T) {
		userID := "test-user-id"
		signedValue := SignCookie(userID)

		parts := strings.Split(signedValue, ".")
		require.Len(t, parts, 2)
		assert.Equal(t, userID, parts[0])
		assert.NotEmpty(t, parts[1])

		assert.True(t, validateCookie(signedValue))
	})
}

func TestValidateCookie(t *testing.T) {
	t.Run("valid cookie", func(t *testing.T) {
		userID := "test-user-id"
		signedValue := SignCookie(userID)
		assert.True(t, validateCookie(signedValue))
	})

	t.Run("invalid format", func(t *testing.T) {
		assert.False(t, validateCookie("invalid-cookie"))
	})

	t.Run("invalid signature", func(t *testing.T) {
		assert.False(t, validateCookie("user-id.invalid-signature"))
	})
}

func TestGetUserID(t *testing.T) {
	t.Run("valid cookie", func(t *testing.T) {
		userID := "test-user-id"
		signedValue := SignCookie(userID)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: cookieName, Value: signedValue})

		gotUserID, ok := GetUserID(req)
		assert.True(t, ok)
		assert.Equal(t, userID, gotUserID)
	})

	t.Run("no cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		gotUserID, ok := GetUserID(req)
		assert.False(t, ok)
		assert.Empty(t, gotUserID)
	})

	t.Run("invalid cookie format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: cookieName, Value: "invalid-cookie"})

		gotUserID, ok := GetUserID(req)
		assert.False(t, ok)
		assert.Empty(t, gotUserID)
	})
}
