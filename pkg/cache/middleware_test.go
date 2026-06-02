package cache

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// MockCache реализует простой in-memory кэш для тестов
type MockCache struct {
	store map[string][]byte
}

func (m *MockCache) Get(key string) []byte {
	return m.store[key]
}

func (m *MockCache) Set(key string, value []byte, duration time.Duration) {
	m.store[key] = value
}

func TestCacheResponse(t *testing.T) {
	// Инициализируем мок кэша
	mockCache := &MockCache{store: make(map[string][]byte)}
	CacheStore = mockCache // Предполагается что CacheStore экспортирован или доступен для подмены

	tests := []struct {
		name             string
		cacheKey         string
		prepopulateCache bool
		duration         string
		nextHandler      http.Handler
		wantStatus       int
		wantBody         string
		wantCached       bool
		wantHeaders      map[string]string
	}{
		// TODO: fix
		/* {
			name:             "cached response exists",
			cacheKey:         "/test",
			prepopulateCache: true,
			duration:         "10s",
			nextHandler:      http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			wantStatus:       http.StatusOK,
			wantBody:         "cached content",
			wantCached:       true,
			wantHeaders:      map[string]string{"X-Test": "123"},
		},
		{
			name:             "uncached response",
			cacheKey:         "/uncached",
			prepopulateCache: false,
			duration:         "5s",
			nextHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-New-Header", "abc")
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte("new content"))
			}),
			wantStatus:  http.StatusCreated,
			wantBody:    "new content",
			wantCached:  true,
			wantHeaders: map[string]string{"X-New-Header": "abc"},
		}, */
		{
			name:             "invalid duration",
			cacheKey:         "/invalid",
			prepopulateCache: false,
			duration:         "invalid",
			nextHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("content"))
			}),
			wantStatus: http.StatusOK,
			wantBody:   "content",
			wantCached: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// Подготовка кэша
			if tt.prepopulateCache {
				CacheStore.Set(tt.cacheKey, []byte(tt.wantBody), 0)
			} else {
				delete(mockCache.store, tt.cacheKey)
			}

			cache := &Cache{duration: tt.duration}
			handler := cache.CacheResponse(tt.nextHandler)

			req := httptest.NewRequest("GET", "http://example.com"+tt.cacheKey, nil)
			rec := httptest.NewRecorder()

			// Act
			handler.ServeHTTP(rec, req)

			// Assert
			// Проверка статуса
			if rec.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}

			// Проверка тела ответа
			if body := rec.Body.String(); body != tt.wantBody {
				t.Errorf("expected body %q, got %q", tt.wantBody, body)
			}

			// Проверка заголовков
			for h, v := range tt.wantHeaders {
				if got := rec.Header().Get(h); got != v {
					t.Errorf("header %s: expected %q, got %q", h, v, got)
				}
			}

			// Проверка кэширования
			if tt.wantCached {
				if content := CacheStore.Get(tt.cacheKey); content == nil {
					t.Error("expected content to be cached, but it's missing")
				}
			} else {
				if content := CacheStore.Get(tt.cacheKey); content != nil {
					t.Error("expected content not to be cached, but it's present")
				}
			}
		})
	}
}

// errorResponseWriter эмулирует ошибки записи
type errorResponseWriter struct {
	*httptest.ResponseRecorder
	writeError bool
	failAfter  int
}

func (w *errorResponseWriter) Write(b []byte) (int, error) {
	if w.failAfter <= 0 {
		w.writeError = true
		return 0, errors.New("simulated write error")
	}
	w.failAfter--
	return w.ResponseRecorder.Write(b)
}

func TestCacheResponseWriteError(t *testing.T) {
	// Arrange
	mockCache := &MockCache{store: make(map[string][]byte)}
	CacheStore = mockCache

	cache := &Cache{duration: "10s"}
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test content"))
	})
	handler := cache.CacheResponse(nextHandler)

	req := httptest.NewRequest("GET", "http://example.com/error-test", nil)
	rec := &errorResponseWriter{
		ResponseRecorder: httptest.NewRecorder(),
		failAfter:        0,
	}

	// Act
	handler.ServeHTTP(rec, req)

	// Assert
	if !rec.writeError {
		t.Error("expected write error, but none occurred")
	}
}
