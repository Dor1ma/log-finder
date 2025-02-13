package repository

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/Dor1ma/log-finder/logfinder/internal/models"
	"github.com/Dor1ma/log-finder/logfinder/pkg/utils"
)

type logFileMetadata struct {
	path  string
	start time.Time
	end   time.Time
}

type LogRepository struct {
	logDir     string
	fileIndex  []logFileMetadata
	indexMutex sync.RWMutex
	fileCache  *fileCache
}

func NewLogRepository(logDir string, maxOpenFiles int, fileCacheTTL time.Duration) (*LogRepository, error) {
	repo := &LogRepository{
		logDir:    logDir,
		fileCache: newFileCache(maxOpenFiles, fileCacheTTL),
	}
	err := repo.RefreshMetadata()

	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *LogRepository) RefreshMetadata() error {
	r.indexMutex.Lock()
	defer r.indexMutex.Unlock()

	files, err := os.ReadDir(r.logDir)
	if err != nil {
		return err
	}

	var newIndex []logFileMetadata
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		path := filepath.Join(r.logDir, f.Name())
		start, end, err := utils.GetFileTimeBounds(path)
		if err != nil {
			continue
		}

		newIndex = append(newIndex, logFileMetadata{
			path:  path,
			start: start,
			end:   end,
		})
	}

	sort.Slice(newIndex, func(i, j int) bool {
		return newIndex[i].start.Before(newIndex[j].start)
	})

	r.fileIndex = newIndex
	return nil
}

func (r *LogRepository) FindByTimestamp(ctx context.Context, t time.Time) (string, error) {
	r.indexMutex.RLock()
	defer r.indexMutex.RUnlock()

	for _, meta := range r.fileIndex {
		if (t.Equal(meta.start) || t.After(meta.start)) &&
			(t.Equal(meta.end) || t.Before(meta.end)) {

			data, err := r.fileCache.Get(meta.path)
			if err != nil {
				return "", err
			}

			result, err := utils.BinarySearchInData(data, t)
			if err != nil {
				return "", err
			}

			return result, nil
		}
	}

	return "", models.ErrNotFound
}
