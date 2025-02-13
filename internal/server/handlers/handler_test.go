package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Dor1ma/log-finder/internal/models"
	"github.com/Dor1ma/log-finder/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRepository struct {
	result     string
	err        error
	refreshErr error
}

func (m *mockRepository) RefreshMetadata() error {
	return m.refreshErr
}

func (m *mockRepository) FindByTimestamp(ctx context.Context, t time.Time) (string, error) {
	return m.result, m.err
}

func TestLogHandler_GetLogByTimestamp(t *testing.T) {
	validTime := "2023-01-01T15:04:05.000"

	tests := []struct {
		name           string
		queryParam     string
		repoResult     string
		repoError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "missing timestamp parameter",
			queryParam:     "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "timestamp parameter is required\n",
		},
		{
			name:           "invalid timestamp format",
			queryParam:     "invalid-time",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid timestamp format\n",
		},
		{
			name:           "log not found",
			queryParam:     validTime,
			repoError:      models.ErrNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "log entry not found\n",
		},
		{
			name:           "internal server error",
			queryParam:     validTime,
			repoError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error\n",
		},
		{
			name:           "successful log retrieval",
			queryParam:     validTime,
			repoResult:     "test log message",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"timestamp":"2023-01-01T15:04:05Z","message":"test log message"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockRepository{
				result: tt.repoResult,
				err:    tt.repoError,
			}

			logService := service.NewLogService(mockRepo, time.Minute)

			handler := NewLogHandler(logService)

			req, err := http.NewRequest("GET", "/logs", nil)
			require.NoError(t, err)

			if tt.queryParam != "" {
				q := req.URL.Query()
				q.Add("timestamp", tt.queryParam)
				req.URL.RawQuery = q.Encode()
			}

			rr := httptest.NewRecorder()
			handler.GetLogByTimestamp(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				assert.JSONEq(t, tt.expectedBody, rr.Body.String())
				assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
			} else {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			}
		})
	}
}

func TestNewLogHandler(t *testing.T) {
	mockRepo := &mockRepository{}
	logService := service.NewLogService(mockRepo, time.Minute)
	handler := NewLogHandler(logService)

	assert.NotNil(t, handler)
}
