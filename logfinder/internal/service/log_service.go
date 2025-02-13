package service

import (
	"context"
	"github.com/Dor1ma/log-finder/logfinder/internal/models"
	"time"
)

var ErrNotFound = models.ErrNotFound

type LogService struct {
	repo  models.LogRepository
	cache *TTLCache
}

func NewLogService(repo models.LogRepository, cacheTTL time.Duration) *LogService {
	return &LogService{
		repo:  repo,
		cache: NewTTLCache(cacheTTL),
	}
}

func (uc *LogService) FindLog(ctx context.Context, timestamp time.Time) (string, error) {
	cacheKey := timestamp.Format(time.RFC3339Nano)

	if entry, ok := uc.cache.Get(cacheKey); ok {
		return entry, nil
	}

	result, err := uc.repo.FindByTimestamp(ctx, timestamp)
	if err != nil {
		return "", err
	}

	uc.cache.Set(cacheKey, result)
	return result, nil
}
