package models

import (
	"context"
	"time"
)

type LogRepository interface {
	FindByTimestamp(ctx context.Context, timestamp time.Time) (string, error)
	RefreshMetadata() error
}
