package utils

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/Dor1ma/log-finder/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBinarySearch(t *testing.T) {
	testData := [][]byte{
		[]byte("2023-01-01T00:00:00.000 line1\n"),
		[]byte("2023-01-01T00:00:01.000 line2\n"),
		[]byte("2023-01-01T00:00:02.000 line3\n"),
	}

	t.Run("exact match", func(t *testing.T) {
		target, _ := time.Parse(timeFormat, "2023-01-01T00:00:01.000")
		result, err := BinarySearchInData(bytes.Join(testData, nil), target)
		assert.NoError(t, err)
		assert.Contains(t, result, "line2")
	})

	t.Run("edge cases", func(t *testing.T) {
		first, err := time.Parse(timeFormat, "2023-01-01T00:00:00.000")
		require.NoError(t, err)
		result, err := BinarySearchInData(bytes.Join(testData, nil), first)
		assert.NoError(t, err)
		assert.Contains(t, result, "line1")

		last, err := time.Parse(timeFormat, "2023-01-01T00:00:02.000")
		require.NoError(t, err)
		result, err = BinarySearchInData(bytes.Join(testData, nil), last)
		assert.NoError(t, err)
		assert.Contains(t, result, "line3")
	})

	t.Run("invalid data handling", func(t *testing.T) {
		invalidData := []byte("invalid log line\n")
		_, err := BinarySearchInData(invalidData, time.Now())
		assert.ErrorIs(t, err, models.ErrInvalidFormat)
	})
}

func TestTimeBounds(t *testing.T) {
	tmpFile := createTestFile(t, []string{
		"2023-01-01T00:00:00.000 first",
		"2023-01-01T00:00:01.000 middle",
		"2023-01-01T00:00:02.000 last",
	})
	defer os.Remove(tmpFile)

	start, end, err := GetFileTimeBounds(tmpFile)
	require.NoError(t, err)

	assert.Equal(t, "2023-01-01T00:00:00.000", start.Format(timeFormat))
	assert.Equal(t, "2023-01-01T00:00:02.000", end.Format(timeFormat))
}

func createTestFile(t *testing.T, lines []string) string {
	f, err := os.CreateTemp("", "test*.log")
	require.NoError(t, err)
	defer f.Close()

	for _, line := range lines {
		_, err = f.WriteString(line + "\n")
		require.NoError(t, err)
	}
	return f.Name()
}
