package limiter

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

// Вспомогательная функция для однократного прохода очистки (для тестирования).
func cleanupVisitorsOnce(timeout time.Duration) {

	mu.Lock()
	defer mu.Unlock()

	for ip, v := range visitors {
		if time.Since(v.lastSeen) > timeout {
			delete(visitors, ip)
		}
	}
}

func TestGetVisitor(t *testing.T) {
	// Сброс глобальной мапы для независимости теста.
	mu.Lock()
	visitors = make(map[string]*visitor)
	mu.Unlock()

	ip := "192.168.0.1"

	limiter1 := getVisitor(ip)
	if limiter1 == nil {
		t.Fatal("getVisitor вернул nil limiter")
	}

	// Сохраняем время первого запроса.
	mu.Lock()
	firstSeen := visitors[ip].lastSeen
	mu.Unlock()

	// Немного подождем и вызовем getVisitor повторно.
	time.Sleep(10 * time.Millisecond)
	limiter2 := getVisitor(ip)
	if limiter2 != limiter1 {
		t.Fatal("Ожидался тот же limiter для одного IP")
	}

	mu.Lock()
	secondSeen := visitors[ip].lastSeen
	mu.Unlock()

	if !secondSeen.After(firstSeen) {
		t.Fatal("lastSeen не обновился при повторном вызове getVisitor")
	}
}

func TestMultipleIPs(t *testing.T) {
	// Сброс глобальной мапы.
	mu.Lock()
	visitors = make(map[string]*visitor)
	mu.Unlock()

	ip1 := "192.168.0.1"
	ip2 := "10.0.0.2"

	limiter1 := getVisitor(ip1)
	limiter2 := getVisitor(ip2)

	if limiter1 == nil || limiter2 == nil {
		t.Fatal("getVisitor вернул nil limiter для одного из IP")
	}

	if limiter1 == limiter2 {
		t.Fatal("Ожидалось, что для разных IP будут созданы разные limiter'ы")
	}

	// Проверяем, что в мапе присутствуют обе записи.
	mu.Lock()
	if len(visitors) != 2 {
		t.Fatalf("Ожидалось 2 записи в visitors, получено %d", len(visitors))
	}
	mu.Unlock()
}

func TestLimitMiddleware(t *testing.T) {
	// Сброс глобальной мапы для теста.
	mu.Lock()
	visitors = make(map[string]*visitor)
	mu.Unlock()

	// Тестовый обработчик, возвращающий статус 200.
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	limitedHandler := Limit(testHandler)

	// Формируем запрос с валидным RemoteAddr.
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = net.JoinHostPort("127.0.0.1", "12345")
	rr := httptest.NewRecorder()

	// Выполняем запрос и проверяем, что он проходит (до исчерпания лимита).
	limitedHandler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Ожидался статус %d, получен %d", http.StatusOK, rr.Code)
	}

	// Теперь проверим превышение лимита.
	// Для лимитера установлено 2 запроса в секунду и burst=5.
	var lastCode int
	requests := 7
	for i := 0; i < requests; i++ {
		limitedHandler.ServeHTTP(rr, req)
		lastCode = rr.Code
		rr = httptest.NewRecorder()
	}
	if lastCode != http.StatusTooManyRequests {
		t.Fatalf("Ожидался статус %d при превышении лимита, получен %d", http.StatusTooManyRequests, lastCode)
	}
}

func TestLimitMiddlewareInvalidRemoteAddr(t *testing.T) {
	// Тестируем ветку, когда RemoteAddr имеет некорректный формат.
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	limitedHandler := Limit(testHandler)

	// Передаем некорректное значение RemoteAddr.
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "invalid_remote_addr"
	rr := httptest.NewRecorder()

	limitedHandler.ServeHTTP(rr, req)
	// Ожидаем ошибку 500.
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("Ожидался статус %d для некорректного RemoteAddr, получен %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestCleanupVisitors(t *testing.T) {
	// Сброс глобальной мапы.
	mu.Lock()
	visitors = make(map[string]*visitor)
	mu.Unlock()

	ip := "10.0.0.1"
	// Добавляем посетителя с устаревшим временем (например, 5 минут назад).
	mu.Lock()
	visitors[ip] = &visitor{
		limiter:  rate.NewLimiter(rate.Limit(1), 1),
		lastSeen: time.Now().Add(-5 * time.Minute),
	}
	mu.Unlock()

	// Вызываем один проход очистки с timeout = 3 минуты.
	cleanupVisitorsOnce(3 * time.Minute)

	mu.Lock()
	_, exists := visitors[ip]
	mu.Unlock()

	if exists {
		t.Fatalf("Ожидалось, что посетитель с IP %s будет удалён", ip)
	}
}
