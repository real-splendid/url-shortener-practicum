package storage

import (
	"testing"

	"github.com/real-splendid/url-shortener-practicum/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorage(t *testing.T) {
	t.Run("create new memory storage", func(t *testing.T) {
		storage := NewMemoryStorage()
		assert.NotNil(t, storage)
		assert.NotNil(t, storage.data)
		assert.NotNil(t, storage.userURLs)
		assert.NotNil(t, storage.deleted)
	})

	t.Run("close storage", func(t *testing.T) {
		storage := NewMemoryStorage()
		err := storage.Close()
		assert.NoError(t, err)
	})

	t.Run("set and get value", func(t *testing.T) {
		storage := NewMemoryStorage()
		userID := "test-user"
		key := "test-key"
		value := "https://example.com"

		_, err := storage.Set(key, value, userID)
		require.NoError(t, err)

		gotValue, err := storage.Get(key)
		require.NoError(t, err)
		assert.Equal(t, value, gotValue)

		urls, err := storage.GetUserURLs(userID)
		require.NoError(t, err)
		require.Len(t, urls, 1)
		assert.Equal(t, key, urls[0].ShortURL)
		assert.Equal(t, value, urls[0].OriginalURL)
	})

	t.Run("get non-existent key", func(t *testing.T) {
		storage := NewMemoryStorage()
		_, err := storage.Get("non-existent-key")
		assert.Error(t, err)
		assert.Equal(t, "key not found", err.Error())
	})

	t.Run("delete user URLs", func(t *testing.T) {
		storage := NewMemoryStorage()
		userID := "test-user"
		key1 := "test-key-1"
		key2 := "test-key-2"
		value1 := "https://example1.com"
		value2 := "https://example2.com"

		_, err := storage.Set(key1, value1, userID)
		require.NoError(t, err)
		_, err = storage.Set(key2, value2, userID)
		require.NoError(t, err)

		err = storage.DeleteUserURLs(userID, []string{key1})
		require.NoError(t, err)

		_, err = storage.Get(key1)
		assert.Equal(t, internal.ErrURLDeleted, err)

		gotValue, err := storage.Get(key2)
		require.NoError(t, err)
		assert.Equal(t, value2, gotValue)

		urls, err := storage.GetUserURLs(userID)
		require.NoError(t, err)
		require.Len(t, urls, 1)
		assert.Equal(t, key2, urls[0].ShortURL)
		assert.Equal(t, value2, urls[0].OriginalURL)
	})

	t.Run("get user URLs for user with no URLs", func(t *testing.T) {
		storage := NewMemoryStorage()
		urls, err := storage.GetUserURLs("non-existent-user")
		require.NoError(t, err)
		assert.Empty(t, urls)
	})

	t.Run("delete all user URLs", func(t *testing.T) {
		storage := NewMemoryStorage()
		userID := "test-user"
		key1 := "test-key-1"
		key2 := "test-key-2"
		value1 := "https://example1.com"
		value2 := "https://example2.com"

		_, err := storage.Set(key1, value1, userID)
		require.NoError(t, err)
		_, err = storage.Set(key2, value2, userID)
		require.NoError(t, err)

		err = storage.DeleteUserURLs(userID, []string{key1, key2})
		require.NoError(t, err)

		urls, err := storage.GetUserURLs(userID)
		require.NoError(t, err)
		assert.Empty(t, urls)
	})
}
