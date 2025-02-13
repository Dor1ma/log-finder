package repository

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Dor1ma/log-finder/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var timeFormat = "2006-01-02T15:04:05.000"

func TestLogRepository(t *testing.T) {
	tmpDir := t.TempDir()
	createTestLogFile(t, tmpDir, "test1.log", []string{
		"2023-01-01T00:00:00.000 line1",
		"2023-01-01T00:00:01.000 line2",
		"2023-01-01T00:00:02.000 line3",
	})

	ctx := context.Background()
	testTime, _ := time.Parse(timeFormat, "2023-01-01T00:00:01.000")

	t.Run("successful creation and basic operations", func(t *testing.T) {
		repo, err := NewLogRepository(tmpDir, 10, time.Minute, time.Hour)
		require.NoError(t, err)
		defer repo.Close()

		assert.Equal(t, 1, repo.FileCount())

		result, err := repo.FindByTimestamp(ctx, testTime)
		require.NoError(t, err)
		assert.Contains(t, result, "line2")
	})

	t.Run("file rotation handling", func(t *testing.T) {
		repo, _ := NewLogRepository(tmpDir, 10, time.Minute, time.Hour)
		defer repo.Close()

		createTestLogFile(t, tmpDir, "test2.log", []string{
			"2023-01-01T00:00:03.000 line4",
		})

		require.NoError(t, repo.RefreshMetadata())
		assert.Equal(t, 2, repo.FileCount())
	})

	t.Run("not found handling", func(t *testing.T) {
		repo, _ := NewLogRepository(tmpDir, 10, time.Minute, time.Hour)
		defer repo.Close()

		invalidTime, _ := time.Parse(timeFormat, "2024-01-01T00:00:00.000")
		_, err := repo.FindByTimestamp(ctx, invalidTime)
		assert.ErrorIs(t, err, models.ErrNotFound)
	})
}

func createTestLogFile(t *testing.T, dir, name string, lines []string) {
	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	require.NoError(t, err)
	defer f.Close()

	for _, line := range lines {
		_, err = f.WriteString(line + "\n")
		require.NoError(t, err)
	}
}
