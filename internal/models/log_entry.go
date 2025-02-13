package models

import "time"

type LogEntry struct {
	Timestamp time.Time
	Message   string
}

func NewLogEntry(timestamp time.Time, message string) *LogEntry {
	return &LogEntry{
		Timestamp: timestamp,
		Message:   message,
	}
}
