package service

import (
	"context"
	"time"

	"github.com/Dor1ma/log-finder/internal/models"
)

var ErrNotFound = models.ErrNotFound
var timeFormat = "2006-01-02T15:04:05.000"

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

func (service *LogService) FindLog(ctx context.Context, timestamp time.Time) (string, error) {
	cacheKey := timestamp.Format(timeFormat)

	if entry, ok := service.cache.Get(cacheKey); ok {
		return entry, nil
	}

	result, err := service.repo.FindByTimestamp(ctx, timestamp)
	if err != nil {
		return "", err
	}

	service.cache.Set(cacheKey, result)
	return result, nil
}
