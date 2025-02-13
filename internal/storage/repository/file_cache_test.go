package repository

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileCache(t *testing.T) {
	t.Run("basic cache operations", func(t *testing.T) {
		tmpDir := t.TempDir()
		lines := []string{
			"2023-01-01T00:00:00.000 line1",
			"2023-01-01T00:00:01.000 line2",
		}
		filePath := createTestLogFileForCache(t, tmpDir, "test_cache.log", lines)

		cache := NewFileCache(2, time.Minute)

		data, err := cache.Get(filePath)
		require.NoError(t, err)
		assert.Contains(t, string(data), "line1")
		assert.Contains(t, string(data), "line2")

		cachedData, err := cache.Get(filePath)
		assert.NoError(t, err)
		assert.Equal(t, data, cachedData, "Cached data should match original")
	})

	t.Run("TTL expiration", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := createTestLogFileForCache(t, tmpDir, "test_ttl.log", []string{"2023-01-01T00:00:00.000 line1"})

		cache := NewFileCache(2, time.Microsecond)
		_, err := cache.Get(filePath)
		require.NoError(t, err)

		time.Sleep(time.Millisecond)

		_, err = cache.Get(filePath)
		assert.ErrorIs(t, err, os.ErrNotExist, "Cache entry should expire")
	})
}

func createTestLogFileForCache(t *testing.T, dir, name string, lines []string) string {
	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	require.NoError(t, err)
	defer f.Close()

	for _, line := range lines {
		_, err = f.WriteString(line + "\n")
		require.NoError(t, err)
	}
	return path
}
