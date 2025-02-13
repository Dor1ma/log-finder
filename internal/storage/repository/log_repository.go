package repository

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/Dor1ma/log-finder/internal/models"
	"github.com/Dor1ma/log-finder/pkg/utils"
)

type logFileMetadata struct {
	path  string
	start time.Time
	end   time.Time
}

type LogRepository struct {
	logDir          string
	fileIndex       []logFileMetadata
	indexMutex      sync.RWMutex
	fileCache       *fileCache
	refreshInterval time.Duration
	done            chan struct{}
	wg              sync.WaitGroup
}

func NewLogRepository(logDir string, maxOpenFiles int, fileCacheTTL, refreshInterval time.Duration) (*LogRepository, error) {
	repo := &LogRepository{
		logDir:          logDir,
		fileCache:       NewFileCache(maxOpenFiles, fileCacheTTL),
		refreshInterval: refreshInterval,
		done:            make(chan struct{}),
	}

	if err := repo.RefreshMetadata(); err != nil {
		return nil, err
	}

	repo.startPeriodicRefresh()
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
			log.Printf("Skipping file %s: %v", path, err)
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
	log.Printf("Metadata refreshed. Files in index: %d", len(r.fileIndex))
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
				return "", models.ErrNotFound
			}

			return result, nil
		}
	}

	return "", models.ErrNotFound
}

func (r *LogRepository) startPeriodicRefresh() {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		ticker := time.NewTicker(r.refreshInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := r.RefreshMetadata(); err != nil {
					log.Printf("Metadata refresh error: %v", err)
				}
			case <-r.done:
				return
			}
		}
	}()
}

func (r *LogRepository) FileCount() int {
	r.indexMutex.RLock()
	defer r.indexMutex.RUnlock()
	return len(r.fileIndex)
}

func (r *LogRepository) Close() {
	close(r.done)
	r.wg.Wait()
	r.fileCache.Clear()
}

func (c *fileCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for k, entry := range c.cache {
		delete(c.cache, k)
		c.lruList.Remove(entry.element)
	}

	c.cache = make(map[string]*cacheEntry)
	c.lruList.Init()
}
