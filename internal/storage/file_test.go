package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/real-splendid/url-shortener-practicum/internal"
	"github.com/stretchr/testify/assert"
)

func createTestFileStorage(t *testing.T, path string) *fileStorage {
	storage, err := NewFileStorage(path)
	if err != nil {
		t.Fatal(err)
	}
	return storage
}

func addTestURL(t *testing.T, storage *fileStorage, key, url, userID string) {
	_, err := storage.Set(key, url, userID)
	if err != nil {
		t.Fatal(err)
	}

	value, err := storage.Get(key)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, url, value)
}

func TestFileStorage(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "file_storage_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Fatal(err)
		}
	}()

	t.Run("create new file storage", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "test_storage.json")

		storage, err := NewFileStorage(filePath)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			err = storage.Close()
			if err != nil {
				t.Fatal(err)
			}
		}()

		assert.NotNil(t, storage)
		assert.NotNil(t, storage.file)
		assert.NotNil(t, storage.records)
		assert.NotNil(t, storage.userURLs)

		_, err = os.Stat(filePath)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("set and get value", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "test_set_get.json")

		storage, err := NewFileStorage(filePath)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			err = storage.Close()
			if err != nil {
				t.Fatal(err)
			}
		}()

		userID := "test-user"
		key := "test-key"
		value := "https://example.com"

		_, err = storage.Set(key, value, userID)
		if err != nil {
			t.Fatal(err)
		}

		gotValue, err := storage.Get(key)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, value, gotValue)

		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatal(err)
		}

		var record fileRecord
		err = json.Unmarshal(fileContent[:len(fileContent)-1], &record)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, key, record.UUID)
		assert.Equal(t, value, record.OriginalURL)
		assert.Equal(t, userID, record.UserID)
		assert.False(t, record.IsDeleted)
	})

	t.Run("get non-existent key", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "test_non_existent.json")

		storage, err := NewFileStorage(filePath)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			err = storage.Close()
			if err != nil {
				t.Fatal(err)
			}
		}()

		_, err = storage.Get("non-existent-key")
		if err == nil {
			t.Fatal("expected error for non-existent key")
		}
		assert.Equal(t, "key not found", err.Error())
	})

	t.Run("duplicate key error", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "test_duplicate.json")

		storage, err := NewFileStorage(filePath)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			err = storage.Close()
			if err != nil {
				t.Fatal(err)
			}
		}()

		userID := "test-user"
		key1 := "test-key-1"
		value := "https://example.com"

		_, err = storage.Set(key1, value, userID)
		if err != nil {
			t.Fatal(err)
		}

		key2 := "test-key-2"
		existingKey, err := storage.Set(key2, value, userID)
		if err == nil {
			t.Fatal("expected error for duplicate key")
		}
		assert.Equal(t, internal.ErrDuplicateKey, err)
		assert.Equal(t, key1, existingKey)
	})

	t.Run("get user URLs", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "test_get_user_urls.json")

		f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			t.Fatal(err)
		}
		err = f.Close()
		if err != nil {
			t.Fatal(err)
		}

		userID := "test-user"
		key1 := "test-key-1"
		value1 := "https://example1.com"

		record := fileRecord{
			UUID:        key1,
			OriginalURL: value1,
			UserID:      userID,
			IsDeleted:   false,
		}

		jsonData, err := json.Marshal(record)
		if err != nil {
			t.Fatal(err)
		}

		err = os.WriteFile(filePath, append(jsonData, '\n'), 0666)
		if err != nil {
			t.Fatal(err)
		}

		storage, err := NewFileStorage(filePath)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			err = storage.Close()
			if err != nil {
				t.Fatal(err)
			}
		}()

		gotValue, err := storage.Get(key1)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, value1, gotValue)

		urls, err := storage.GetUserURLs(userID)
		if err != nil {
			t.Fatal(err)
		}
		assert.Len(t, urls, 1)
		assert.Equal(t, value1, urls[0].OriginalURL)
		assert.Equal(t, key1, urls[0].ShortURL)

		urls, err = storage.GetUserURLs("non-existent-user")
		if err != nil {
			t.Fatal(err)
		}
		assert.Empty(t, urls)

		err = storage.DeleteUserURLs(userID, []string{key1})
		if err != nil {
			t.Fatal(err)
		}
		urls, err = storage.GetUserURLs(userID)
		if err != nil {
			t.Fatal(err)
		}
		assert.Empty(t, urls)
	})

	t.Run("delete user URLs", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "test_delete_user_urls.json")

		storage, err := NewFileStorage(filePath)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			err = storage.Close()
			if err != nil {
				t.Fatal(err)
			}
		}()

		userID := "test-user"
		key1 := "test-key-1"
		value1 := "https://example.com/1"
		key2 := "test-key-2"
		value2 := "https://example.com/2"

		addTestURL(t, storage, key1, value1, userID)
		addTestURL(t, storage, key2, value2, userID)

		err = storage.DeleteUserURLs(userID, []string{key1})
		if err != nil {
			t.Fatal(err)
		}

		_, err = storage.Get(key1)
		if err == nil {
			t.Fatal("expected error for deleted URL")
		}
		assert.Equal(t, "url deleted", err.Error())

		gotValue, err := storage.Get(key2)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, value2, gotValue)
	})

	t.Run("delete non-existent URL", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "test_delete_non_existent.json")

		f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			t.Fatal(err)
		}
		err = f.Close()
		if err != nil {
			t.Fatal(err)
		}

		storage, err := NewFileStorage(filePath)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			err = storage.Close()
			if err != nil {
				t.Fatal(err)
			}
		}()

		err = storage.DeleteUserURLs("test-user", []string{"non-existent-key"})
		if err == nil {
			t.Fatal("expected error for non-existent key")
		}
		assert.Equal(t, "record not found by short key", err.Error())
	})

	t.Run("load existing data", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "test_load_existing.json")

		storage1, err := NewFileStorage(filePath)
		if err != nil {
			t.Fatal(err)
		}

		userID := "test-user"
		key := "test-key"
		value := "https://example.com"

		_, err = storage1.Set(key, value, userID)
		if err != nil {
			t.Fatal(err)
		}

		err = storage1.Close()
		if err != nil {
			t.Fatal(err)
		}

		storage2, err := NewFileStorage(filePath)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			err = storage2.Close()
			if err != nil {
				t.Fatal(err)
			}
		}()

		gotValue, err := storage2.Get(key)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, value, gotValue)
	})

	t.Run("get key by value", func(t *testing.T) {
		storage, err := NewFileStorage("testfile.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			err = os.Remove("testfile.txt")
			if err != nil {
				t.Fatal(err)
			}
			err = storage.Close()
			if err != nil {
				t.Fatal(err)
			}
		}()

		key := "testKey"
		value := "testValue"
		userID := "testUser"
		_, err = storage.Set(key, value, userID)
		if err != nil {
			t.Fatal(err)
		}

		foundKey, err := storage.getKeyByValue(value)
		if err != nil {
			t.Fatal(err)
		}
		if foundKey != key {
			t.Errorf("Expected key %s, got %s", key, foundKey)
		}

		_, err = storage.getKeyByValue("nonExistentValue")
		if err == nil {
			t.Fatal("expected error for non-existent value")
		}
	})
}

func TestFileStorage_Close(t *testing.T) {
	storage, err := NewFileStorage("testfile.txt")
	if err != nil {
		t.Fatal(err)
	}

	err = storage.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = storage.Close()
	if err == nil {
		t.Fatal("expected error for already closed file")
	}
	assert.Equal(t, "close testfile.txt: file already closed", err.Error())

	err = os.Remove("testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewFileStorage_Error(t *testing.T) {
	_, err := NewFileStorage("/invalid/path/testfile.txt")
	if err == nil {
		t.Fatal("expected error for invalid file path")
	}
}

func TestFileStorage_DeleteUserURLs(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		storage, err := NewFileStorage("testfile.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			err = os.Remove("testfile.txt")
			if err != nil {
				t.Fatal(err)
			}
			err = storage.Close()
			if err != nil {
				t.Fatal(err)
			}
		}()

		keys := []string{"key1", "key2"}
		urls := []string{"http://example.com/1", "http://example.com/2"}
		userID := "testUser"
		for i, key := range keys {
			_, err = storage.Set(key, urls[i], userID)
			if err != nil {
				t.Fatal(err)
			}
		}

		err = storage.DeleteUserURLs(userID, keys)
		if err != nil {
			t.Fatal(err)
		}

		for _, key := range keys {
			_, err = storage.Get(key)
			if !errors.Is(err, internal.ErrURLDeleted) {
				t.Errorf("Expected URL to be deleted, got error: %v", err)
			}
		}
	})

	t.Run("non-existent key", func(t *testing.T) {
		storage, err := NewFileStorage("testfile.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			err = os.Remove("testfile.txt")
			if err != nil {
				t.Fatal(err)
			}
			err = storage.Close()
			if err != nil {
				t.Fatal(err)
			}
		}()

		err = storage.DeleteUserURLs("testUser", []string{"nonExistentKey"})
		if err == nil {
			t.Fatal("expected error for non-existent key")
		}
	})

	t.Run("empty keys list", func(t *testing.T) {
		storage, err := NewFileStorage("testfile.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			err = os.Remove("testfile.txt")
			if err != nil {
				t.Fatal(err)
			}
			err = storage.Close()
			if err != nil {
				t.Fatal(err)
			}
		}()

		err = storage.DeleteUserURLs("testUser", []string{})
		if err != nil {
			t.Fatal(err)
		}
	})
}
