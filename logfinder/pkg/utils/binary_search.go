package utils

import (
	"bufio"
	"bytes"
	"os"
	"time"

	"github.com/Dor1ma/log-finder/logfinder/internal/models"
)

var timeFormat = "2006-01-02T15:04:05.000"

func GetFileTimeBounds(path string) (time.Time, time.Time, error) {
	file, err := os.Open(path)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var firstLine string
	if scanner.Scan() {
		firstLine = scanner.Text()
	}

	var lastLine string
	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	if lastLine == "" {
		lastLine = firstLine
	}

	start, err := ParseTimestamp(firstLine)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	end, err := ParseTimestamp(lastLine)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return start, end, nil
}

func ParseTimestamp(line string) (time.Time, error) {
	if len(line) < 23 {
		return time.Time{}, models.ErrInvalidFormat
	}
	return time.Parse(timeFormat, line[:23])
}

func BinarySearchInData(data []byte, target time.Time) (string, error) {
	lines := bytes.Split(data, []byte{'\n'})
	low := 0
	high := len(lines) - 1

	for low <= high {
		mid := (low + high) / 2
		line := string(lines[mid])

		lineTime, err := ParseTimestamp(line)
		if err != nil {
			return "", err
		}

		if lineTime.Equal(target) {
			return line, nil
		}

		if target.Before(lineTime) {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}

	return "", models.ErrNotFound
}
