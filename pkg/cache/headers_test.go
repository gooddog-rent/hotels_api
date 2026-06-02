package cache

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockHandler реализует http.Handler для проверки вызова следующего обработчика
type mockHandler struct {
	called bool
}

func (m *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.called = true
	w.WriteHeader(http.StatusOK)
}

func TestCacheHeaders(t *testing.T) {

	// Arrange
	tests := []struct {
		name           string
		duration       string
		wantStatus     int
		wantHeader     string
		wantBody       string
		wantNextCalled bool
	}{
		{
			name:           "valid duration",
			duration:       "60s",
			wantStatus:     http.StatusOK,
			wantHeader:     "private, max-age=60",
			wantBody:       "",
			wantNextCalled: true,
		},
		{
			name:           "invalid duration",
			duration:       "invalid",
			wantStatus:     http.StatusInternalServerError,
			wantHeader:     "",
			wantBody:       "Internal Server Error\n",
			wantNextCalled: false,
		},
		{
			name:           "zero duration",
			duration:       "0s",
			wantStatus:     http.StatusOK,
			wantHeader:     "private, max-age=0",
			wantBody:       "",
			wantNextCalled: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			// Arrange
			mockNext := &mockHandler{}
			cache := Cache{duration: tt.duration}
			handler := cache.CacheHeaders(mockNext)

			req := httptest.NewRequest("GET", "http://example.com", nil)
			rec := httptest.NewRecorder()

			// Act
			handler.ServeHTTP(rec, req)

			// Assert
			// Проверяем статус код
			if rec.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}

			// Проверяем заголовок Cache-Control
			gotHeader := rec.Header().Get("Cache-Control")
			if gotHeader != tt.wantHeader {
				t.Errorf("expected Cache-Control header '%s', got '%s'", tt.wantHeader, gotHeader)
			}

			// Проверяем тело ответа
			if rec.Body.String() != tt.wantBody {
				t.Errorf("expected body '%s', got '%s'", tt.wantBody, rec.Body.String())
			}

			// Проверяем вызов следующего обработчика
			if mockNext.called != tt.wantNextCalled {
				t.Errorf("expected next called %t, got %t", tt.wantNextCalled, mockNext.called)
			}
		})
	}
}

func TestCacheHeadersWithNextHandlerModifications(t *testing.T) {

	// Arrange
	// Проверка что middleware не перезаписывает существующие заголовки
	mockNext := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)
	})

	cache := Cache{duration: "30s"}
	handler := cache.CacheHeaders(mockNext)

	req := httptest.NewRequest("GET", "http://example.com", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(rec, req)

	// Assert
	// Проверяем что заголовок от следующего обработчика имеет приоритет
	if got := rec.Header().Get("Cache-Control"); got != "no-cache" {
		t.Errorf("expected Cache-Control header from next handler, got '%s'", got)
	}
}
